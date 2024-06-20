var lat = parseFloat(document.getElementById("lat").dataset.name);
var long = parseFloat(document.getElementById("long").dataset.name);
var china = document.getElementById("china").dataset.name;
var map;
if (china == "true") {
  var getTiles = function (baseUrl, tilePoint) {
    String.prototype.formatString = function () {
      let a = this;
      let b;
      for (b in arguments) {
        a = a.replace(/{[0-9]}/, arguments[b]);
      }
      return a;
    };

    tilePoint = baseUrl;
    baseUrl = this._url;
    var offset = Math.pow(2, tilePoint.z - 1);
    var x = tilePoint.x - offset;
    var y = offset - tilePoint.y - 1;
    var z = tilePoint.z;
    return baseUrl.formatString(parseInt(Math.random() * 10) % 4, x, y, z);
  };

  var subdomains = "0123";
  var baiduCRS = new L.Proj.CRS(
    "EPSG:3857",
    "+proj=merc +a=6378206 +b=6356584.314245179 +lat_ts=0.0 +lon_0=0.0 +x_0=0 +y_0=0 +k=1.0 +units=m +nadgrids=@null +wktext  +no_defs",
    {
      resolutions: (function () {
        var res = [];
        for (var i = 0; i < 20; ++i) {
          res[i] = Math.pow(2, 18 - i);
        }
        return res;
      })(),
      origin: [-3.3554432e7, 3.3554432e7],
      bounds: L.bounds(
        [-3.3554432e7, 3.3554432e7],
        [3.3554432e7, -3.3554432e7],
      ),
    },
  );
  var baiduLayer = L.tileLayer(
    "http://maponline{0}.bdimg.com/tile/?qt=tile&x={1}&y={2}&z={3}&styles=pl",
    {
      minZoom: 3,
      noWrap: false,
    },
  );

  baiduLayer.getTileUrl = getTiles;
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
