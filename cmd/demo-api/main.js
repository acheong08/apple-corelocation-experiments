var lat = parseFloat(document.getElementById("lat").dataset.name);
var long = parseFloat(document.getElementById("long").dataset.name);
var china = document.getElementById("china").dataset.name;
var map;
if (china) {
  var subdomains = "0123";
  var baiduCRS = new L.Proj.CRS(
    "EPSG:900913",
    "+proj=merc +a=6378206 +b=6356584.314245179 +lat_ts=0.0 +lon_0=0.0 +x_0=0 +y_0=0 +k=1.0 +units=m +nadgrids=@null +wktext  +no_defs",
    {
      resolutions: (function () {
        level = 19;
        var res = [];
        res[0] = Math.pow(2, 18);
        for (var i = 1; i < level; i++) {
          res[i] = Math.pow(2, 18 - i);
        }
        return res;
      })(),
      origin: [0, 0],
      bounds: L.bounds([20037508.342789244, 0], [0, 20037508.342789244]),
    },
  );
  var baiduLayer = L.tileLayer(
    "http://maponline{s}.bdimg.com/tile/?qt=tile&x={x}&y={y}&z={z}&styles=pl&scaler=1&p=1",
    {
      name: "vec",
      subdomains: subdomains,
      tms: true,
    },
  );
  map = L.map("map", {
    crs: baiduCRS,
  })
    .setView([lat, long], 13)
    .addLayer(baiduLayer);
} else {
  var osmLayer = L.tileLayer("https://tile.openstreetmap.org/{z}/{x}/{y}.png", {
    maxZoom: 19,
    attribution:
      '&copy; <a href="http://www.openstreetmap.org/copyright">OpenStreetMap</a>',
  });
  map = L.map("map", {}).setView([lat, long], 13).addLayer(osmLayer);
}

function plotPoint(p) {
  L.marker([p.y, p.x]).addTo(map).bindPopup(p.id);
}
let blocker = false;
async function onMapClick(e) {
  if (blocker) return;
  blocker = true;
  // Create a loading dialog
  const loading = document.createElement("dialog");
  loading.innerHTML = "Loading...";
  document.body.appendChild(loading);
  loading.showModal();

  const cleanUp = () => {
    blocker = false;
    loading.remove();
  };

  const lat = e.latlng.lat;
  const long = e.latlng.lng;
  let resp = await fetch("/gps", {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify({ lat, long }),
  });
  // Check status code
  if (!resp.ok) {
    cleanUp();
    alert("Error: " + resp.status);
    return;
  }
  resp = await resp.json();
  console.log(resp);
  plotPoint(resp.closest);
  resp.points.forEach(plotPoint);
  cleanUp();
}
map.on("click", onMapClick);
