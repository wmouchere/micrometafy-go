"use strict"
let serverPlaylist = "http://localhost:8081/micrometafy-playlist/api/";

var playlistTracks;
var currentPlaylist = -1;
var playlistIndex = 0;

var audio = document.getElementById("audioplayer");

document.getElementById('query').onkeydown = function(e){
    if(e.keyCode == 13){
        submit()
    }
};

audio.addEventListener('ended',function(){
    if(currentPlaylist != -1) {
        let urls = playlistTracks.get(currentPlaylist)
        if(playlistIndex < urls.length) {
            let info = `<source src=${urls[playlistIndex]} type="audio/mpeg">`
            audio.innerHTML = info
            playlistIndex++;
            audio.load();
            audio.play();
        } else {
            audio.innerHTML = "";
            audio.load();
            currentPlaylist = -1;
            playlistIndex = 0;
        }
    }
});

function submit(){
    let query = document.getElementById('query').value;
    let req = new XMLHttpRequest();

    req.onreadystatechange = function () {
        if (req.readyState !== 4) return;

        if (req.status === 200) {
            console.log(req.responseType + " " + req.response);
            findAllPlaylists()
        }
    };

    req.open("POST", serverPlaylist + "playlist", true);
    req.setRequestHeader("Content-Type", "application/json;charset=UTF-8");
    req.send(JSON.stringify({name: query, tracks: []}));
}

function findAllPlaylists() {
    let req = new XMLHttpRequest();

    playlistTracks = new Map();

    req.onreadystatechange = function () {
        if (req.readyState !== 4) return;

        if (req.status === 200) {
            console.log(req.responseType + " " + req.response);

            let playlists = JSON.parse(req.response);
            console.log(playlists);
            
            let info = "";
            if (playlists !== null) {
                playlists.forEach(function (p) {
                    let urls = [];
                    p.tracks.forEach(function (t) {
                        urls.push(t.url);
                    });
                    playlistTracks.set(p.id, urls);
                    info += `<div name="playlist">`
                    info += `<label class="playlist-title">${p.name}\t</label>`
                    info += `<button type="submit" onclick="playTracks('${p.id}')" class="pure-button"><i class="fas fa-play"></i></button>`
                    info += `<button type="submit" value="Supprimer" class="pure-button button-delete"><i class="fas fa-trash"></i></button>`
                    info += `<input name="data" type="hidden" value=${encodeURIComponent(JSON.stringify(p))} />`
                    info += `<ol>`
                    p.tracks.forEach(function (t) {
                        info += `<li>
                        <img src="/img/${t.origin}.png" class="api-image">
                        ${t.name} -- ${t.author}   
                        <i class="fas fa-clock"></i>
                        ${Math.floor(t.duration/60000)}min${Math.floor(t.duration%60000/1000)}s.\t`

                        //Display play button only if link for preview is available
                        if(t.url != undefined){
                            info += `<button type="submit" onclick="playTrack('${t.url}')" class="pure-button"><i class="fas fa-play"></i></button>`
                        }
                        info += `</li>`;
                    })
                    
                    info += `</ol>`
                    info += `</div>`
                })
            }
            document.getElementById('info').innerHTML = info;
        } else {
            document.getElementById('info').innerHTML = "Cannot be retrieved";
        }

        var playlists = document.getElementsByName('playlist');
        for(let i = 0; i < playlists.length; i++) {
            let data = JSON.parse(decodeURIComponent(playlists[i].getElementsByTagName('input')[0].getAttribute("value")));
            console.log(data);
            playlists[i].getElementsByTagName('button')[1].onclick = function () {
                let req2 = new XMLHttpRequest();
                req2.onreadystatechange = update();
                req2.open("DELETE", serverPlaylist + "playlist/" + data.id, true);
                req2.send();
                findAllPlaylists()
            }
        }
    };
    req.open("GET", serverPlaylist + "playlists", true);
    req.send();
}

var update = function () {
    return function() {
        findAllPlaylists();
    }
};

function playTrack(url) {
    let info = `<source src=${url} type="audio/mpeg">`
    audio.innerHTML = info
    playlistIndex = 0;
    currentPlaylist = -1;
    audio.pause()
    audio.load()
    audio.play()
}

function playTracks(id) {
    let urls = playlistTracks.get(id)
    let info = `<source src=${urls[0]} type="audio/mpeg">`
    currentPlaylist = id;
    playlistIndex = 1;
    audio.innerHTML = info
    audio.pause()
    audio.load()
    audio.play()
}

findAllPlaylists();