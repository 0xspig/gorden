import ForceGraph3D from "3d-force-graph";

const Graph = new ForceGraph3D(document.getElementById('view'))
  .width(document.getElementById("view").getBoundingClientRect().width)
  .height(document.getElementById("view").getBoundingClientRect().height)
  .showNavInfo(false);
var parsedJson;

//colors
//$color-bg-mid: #626d64;
//$color-bg-white: #fbf0e3;
//$color-text: #323738;
//$color-text-light: #f9f1e9;
//$color-text-dark: #3e4642;
//$color-text-link: #bf9b6e;
//$color-text-highlight: #b9c3bb;
Graph.backgroundColor("#626d64");

Graph.onNodeClick(node => {
  targetNode(node.id);
});



var canvas = Graph.renderer().domElement
canvas.id = "scene";

var view = document.getElementById("view");
Graph.width(canvas.getBoundingClientRect().width)
Graph.height(canvas.getBoundingClientRect().height);
Graph.camera().aspect = canvas.clientWidth / canvas.clientHeight;
Graph.camera().updateProjectionMatrix()

function resizeWindow(){
  console.log("resize");
  var real_style = getComputedStyle(document.getElementById("view"));
  document.getElementById("scene").style.width = real_style.width;
  document.getElementById("scene").style.height = real_style.height;
  Graph.width(document.getElementsByClassName("scene-container")[0].getBoundingClientRect().width);
  Graph.height(document.getElementsByClassName("scene-container")[0].getBoundingClientRect().height);
  Graph.width(document.getElementById("scene").getBoundingClientRect().width);
  Graph.height(document.getElementById("scene").getBoundingClientRect().height);
  Graph.camera().aspect = canvas.clientWidth / canvas.clientHeight;
  Graph.camera().updateProjectionMatrix()
}

addEventListener("resize", resizeWindow)

var xmlhttp = new XMLHttpRequest;
xmlhttp.onreadystatechange = function() {
    if (this.readyState == 4 && this.status == 200) {
		Graph.graphData(JSON.parse(this.responseText));
    console.log("graph data updated");
    pushGraphParams();
    // target uri node (defaults to home if not found)
    var path = window.location.pathname.split('/');
    var id = path[path.length - 1];
    targetNode(id);
    }
    resizeWindow();
};

xmlhttp.open("GET", "/graph-json", true);
xmlhttp.send();

function pushGraphParams(){

    Graph.nodeColor(node => {
      if (node.targeted == true){
        return "#ff33ff";
      }
      if (node.highlighted == true){
        return "#eeeeff"
      }
      //md files
      if (node.data_type == 1){
        return "#fff380";
      }
      // tags
      if (node.data_type == 2){
        return "#6b93c6";
      }
      // categories
      if (node.data_type == 3){
        return "#63335c";
      }
      // external
      if (node.data_type == 4){
        return "#68d588";
      }
    });

    Graph.nodeVal(node => {
      if (node.targeted == true){
        return 8;
      }
      if (node.highlighted == true){
        return 8;
      }
      //md files
      if (node.data_type == 1){
        return 2;
      }
      // tags
      if (node.data_type == 2){
        return 1;
      }
      // categories
      if (node.data_type == 3){
        return 4;
      }
      // external
      if (node.data_type == 4){
        return 1;
      }
    })
}

export function highlightNode(nodeID){
    var target_node;

    console.log("highlighting Node: "+nodeID);
    Graph.graphData().nodes.forEach(node => {
      if (node.id == nodeID){
        target_node = node;
        target_node.highlighted = true;
      }else{
        node.highlighted = false;
      }
    });
    pushGraphParams();
}

export function targetNode(nodeID){

    var target_node;

    console.log("targeting Node: "+nodeID);
    Graph.graphData().nodes.forEach(node => {
      if (node.id == nodeID){
        target_node = node;
        // if external
        if (target_node.data_type == 4) {
          window.open(node.source, '_blank').focus();
          return;
        }
        target_node.targeted = true;
      }else{
        node.targeted = false;
      }
    });
    if (target_node == null){
      targetNode("home.md");
      return false;
    }
    console.log(target_node);
    if (target_node.data_type == 4){
      return false;
    } 
    // Aim at target_node from outside it
    const distance = 225;
    const distRatio = 1 + distance/Math.hypot(target_node.x, target_node.y, target_node.z);

    const newPos = target_node.x || target_node.y || target_node.z
      ? { x: target_node.x * distRatio, y: target_node.y * distRatio, z: target_node.z * distRatio }
      : { x: 0, y: 0, z: distance }; // special case if target_node is in (0,0,0)

    Graph.cameraPosition(
      newPos, // new position
      target_node, // lookAt ({ x, y, z })
      1800  // ms transition duration
    );

    // ui stuff
    getNodeData(target_node);
    pushGraphParams();
    window.history.pushState(null, target_node.name, target_node.id);
    return false;
}

addEventListener("popstate", (event) => {
    // target uri node (defaults to home if not found)
    var path = window.location.pathname.split('/');
    var id = path[path.length - 1];
    targetNode(id);
});


/* !!!!! TAB FUNCTION !!!
  I'm removing this for the time being since its ugly and doesn't add anything to the UX

document.getElementById("data-tab").addEventListener("click", dataToggle)
var dataViewing = true;
function dataToggle(){
  if (dataViewing){
    document.getElementById("data-content").style.flexBasis="0%";
    document.getElementById("data-content").style.display ="none";
    var elemChildren = document.getElementById("data-content").children
    for (var i = 0; i < elemChildren.length; i++){
      elemChildren[i].style.display = "none";
    }
    dataViewing=false
  }else{
    document.getElementById("data-content").style.flexBasis="45%";
    document.getElementById("data-content").style.display ="inline";
    var elemChildren = document.getElementById("data-content").children
    for (var i = 0; i < elemChildren.length; i++){
      elemChildren[i].style.display = "inline";
    }
    dataViewing=true
  }
    resizeWindow()
}

*/

function getNodeData(node){
  var resp = document.getElementById("response");
  if (resp != null){
    resp.style = "opacity: 0%;";
  }
  var node_data_request = new XMLHttpRequest;
  node_data_request.onreadystatechange = function() {
    switch (node_data_request.readyState){
      default:
        return;
      case 3:
        return;
      case 4:
        var content = document.getElementById("content");
        content.innerHTML = node_data_request.responseText;
        var resp = document.getElementById("response");

        var margin = parseFloat(getComputedStyle(content).paddingTop) 
                      + parseFloat(getComputedStyle(content).paddingBottom)
                      + parseFloat(getComputedStyle(resp).marginTop)
                      + parseFloat(getComputedStyle(resp).marginBottom);
        content.style = "height: " + (resp.scrollHeight + margin) + "px;";
        resp.style = "opacity: 100%;";

        return;

    }
  }
  node_data_request.open("GET", "/node-data/"+node.id, true);
  node_data_request.send();

  var node_links_request = new XMLHttpRequest;
  node_links_request.onreadystatechange = function() {
    document.getElementById("link-data").innerHTML = node_links_request.responseText;
  }
  node_links_request.open("GET", "/node-links/"+node.id, true)
  node_links_request.send()
}