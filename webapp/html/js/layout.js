let header = `
<div>
  <a href="index.html">Home</a>
  <a href="search.html">Search</a>
  <a href="playlists.html">Playlists</a>
</div>
`;

let footer = `
<p>&copy; Copyright 2018 - Polytechnique Montréal</p>`;

document.getElementById('header').innerHTML = header;
document.getElementById('footer').innerHTML = footer;

// Sets the current page as the active page in the navbar
var linksArray = document.getElementById('header').getElementsByTagName('a');
for(var i = 0; i < linksArray.length; i++) {
    if (linksArray[i].href == location.href){ linksArray[i].className="active"; }
}