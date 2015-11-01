(function(){

var tab=document.createElement("div")
tab.style.height="30px";
tab.innerHTML="show panel";
tab.style.position="fixed";
tab.style.backgroundColor="cornflowerblue";
tab.style.top="0";
tab.style.left="0";
tab.style.zIndex="499";
tab.style.textDecoration="underline";
tab.style.color="ghostwhite";

document.body.appendChild(tab);

var remove=document.createElement("span");
remove.innerHTML="close panel";
remove.style.textDecoration="underline";
remove.style.color="ghostwhite";



var ws = new WebSocket("ws://localhost/question/");
var el = document.createElement("div");

el.appendChild(remove);

el.style.position="fixed";
el.style.width="600px";

el.style.backgroundColor="cornflowerblue";
el.style.zIndex="500";
el.style.padding="6px";
el.style.top="0";
el.style.left="0";
el.style.boxSizing="border-box";

var wrapper = document.createElement("div");
wrapper.style.position="relative";
wrapper.style.boxSizing="border-box";

var innerel = document.createElement("div");
innerel.style.backgroundColor="ghostwhite";
innerel.style.width="100%";
innerel.style.display="block";
innerel.style.height="500px"
innerel.style.overflow="auto";
innerel.style.padding="6px";
innerel.style.boxSizing="border-box";

var inputel = document.createElement("input");
inputel.type="text";
inputel.style.width="100%";
inputel.style.padding="8px";
inputel.style.margin="4px 0";
inputel.style.display="block";
inputel.style.zIndex="505";
inputel.style.boxSizing="border-box";


inputel.onkeydown=function(e){
if (e.which==13){
    e.preventDefault();
    var msg=inputel.value;
    inputel.value="";
    ws.send(msg);
}
};

document.body.appendChild(el);

wrapper.appendChild(innerel);
wrapper.appendChild(inputel);
el.appendChild(wrapper);

function addLine(msg){
    var span = document.createElement("span");
    span.innerHTML=msg;
    span.style.display="block";
    var isatbottom = false;
    
    if (innerel.scrollTop===innerel.scrollHeight-innerel.clientHeight){
        isatbottom=true;
    }
    
    innerel.appendChild(span);
    
    if (isatbottom){
        innerel.scrollTop = innerel.scrollHeight;
    }
}

ws.onmessage=function(e){
    console.log(e);
    addLine(e.data);
};

tab.onclick=function(e){
    document.body.appendChild(el);
};

remove.onclick=function(e){
    document.body.removeChild(el);
};

}())