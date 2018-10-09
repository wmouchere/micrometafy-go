"use strict"
let serverQuery = "http://localhost:8080/micrometafy-query/api/";
let serverPlaylist = "http://localhost:8081/micrometafy-playlist/api/";

let playlists = findPlaylists();
let tracks;

document.getElementById('query').onkeydown = function(e){
    if(e.keyCode == 13){
      submit()
    }
 };

 function submit(){
    let query = document.getElementById('query').value;
    let req = new XMLHttpRequest();

    req.onreadystatechange = function () {
        if (req.readyState !== 4) return;

        if (req.status === 200) {
            console.log(req.responseType + " " + req.response);

            tracks = JSON.parse(req.response);
            console.log(tracks);

            let info = "";
            let i = 0;
            info += `<h1>Results</h1>`
            info += `<ol>`
            tracks.forEach(function (t) {
                info += `<li>
                            <img src="/img/${t.origin}.png" class="api-image">
                            ${t.name} -- ${t.author}   
                            <i class="fas fa-clock"></i>
                            ${Math.floor(t.duration/60000)}min${Math.floor(t.duration%60000/1000)}s.\t`

                //Display play button only if link for preview is available
                if(t.url != undefined){
                    info += `<button type="submit" onclick="playTrack('${t.url}')" class="pure-button"><i class="fas fa-play"></i></button>`
                }
                
                info += `<div class="dropdown">
                            <button type="submit" onclick="dropDown(${i})" class="dropbtn pure-button button-add">
                                <i class="fas fa-plus dropbtn"></i>
                            </button>
                            <div id="dropdown${i}" class="dropdown-content">
                            <a href="playlists.html" class="button"><i class="fas fa-plus"></i>\tNew playlist</a>`
                            
                playlists.forEach(function (p) {
                                info += `<label onclick="add(${i}, '${p.id}')">${p.name}</label>`
                });

                info +=     `</div>
                        </div>`       
                            
                info += `</li>`;
                i++;
            });
            info += `</ol>`
            document.getElementById('info').innerHTML = info;
        } else {
            document.getElementById('info').innerHTML = "Cannot be retrieved";
        }
    };

    req.open("GET", serverQuery + "search/" + query, true);
    req.send();
}

function findPlaylists(){
    let req = new XMLHttpRequest();

    req.onreadystatechange = function () {
        if (req.readyState !== 4) return;

        if (req.status === 200) {
            console.log(req.responseType + " " + req.response);

            playlists = JSON.parse(req.response);
            console.log(playlists);
            return playlists;
        }
    }
    
    req.open("GET", serverPlaylist + "playlists", true);
    req.send();
}

function dropDown(i){
    let id = "dropdown" + i;
    document.getElementById(id).classList.toggle("show");
}

function playTrack(url) {
    let info = `<source src=${url} type="audio/mpeg">`
    let audio = document.getElementById("audioplayer")
    audio.innerHTML = info
    audio.pause()
    audio.load()
    audio.play()
}

function add(t, p){
    let req = new XMLHttpRequest();
    let track = tracks[t];

    req.onreadystatechange = function () {
        if (req.readyState !== 4) return;

        if (req.status === 200) {
            console.log(req.responseType + " " + req.response);
        }
    };

    req.open("PUT", serverPlaylist + "playlist/" + p + "/add", true);
    req.setRequestHeader("Content-Type", "application/json;charset=UTF-8");
    console.log(track);
    console.log(JSON.stringify({track}));
    req.send(JSON.stringify(track));
    

}

window.onclick = function(event) {
    if (!event.target.matches('.dropbtn')) {
  
      var dropdowns = document.getElementsByClassName("dropdown-content");
      var i;
      for (i = 0; i < dropdowns.length; i++) {
        var openDropdown = dropdowns[i];
        if (openDropdown.classList.contains('show')) {
          openDropdown.classList.remove('show');
        }
      }
    }
  }