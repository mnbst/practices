// some let
let player;
let href;
let endtime;
let starttime;
let text;
let repeat = false;
let s_id = -1;
let content;
let timer;
let subClicked = false;
let textClicked = false;
let currentTime=0;
let subTitles;

let tag = document.createElement('script');
tag.src = "http://www.youtube.com/player_api";
let firstScriptTag = document.getElementsByTagName('script')[0];
firstScriptTag.parentNode.insertBefore(tag, firstScriptTag);

$(function(){
    $('.text').click(function () {
    $('.col-md-12 text-center, .col-md-12, #subtitle_detail_area').remove();  
    $(this).before("<div class='col-md-12 text-center'>\
    <div class='embed-responsive embed-responsive-16by9'>\
    <div id='player'></div></div></div><div id='subtitle_area' class='col-md-12'>\
    <div class='panel panel-default'><div id='subtitle_block'></div></div>\
    </div><div id='subtitle_detail_area'></div>");

    $("html,body").animate({scrollTop:$('#player')['0'].parentNode.parentNode.offsetTop-60});
    subTitles=$("#subtitle_block");
    href =  $(this).attr("id");
    id =  this.children[1].id;

    $.ajax({
        url: window.location.pathname,
        contentType: "application/json",   
        type: 'GET',      
        data: {'id': href},
        dataType: 'json',  
        timeout: 10000, 
        cache: false, 
    }).then(
        function(jsondata) {
            $('#json-data').html("");
            $('#json-data').html(JSON.stringify(jsondata));
            let obj =document.getElementById("json-data");
            content = JSON.parse(obj.textContent);
            content.forEach(function (element, index) {
                $("#subtitle_block").append("<tr><td>" + "<a id='" + index + "'' class='subtitle'>" + "<i class='far fa-play-circle fa-lg'></i>" + "</a></td>" + "<td><a id='" + index + "'' class='subtitle'><span id='subtitle_text_" + index + "' >" + element["text"] + "</span></a></td></tr>");
            });     
            caption_detail(id);
            colorSub(id);
            scrollSub(id); 
            player = new YT.Player('player', {
                height: '200',
                width: '320',
                playerVars: {
                    autoplay: 1,
                    modestbranding: 1,
                    rel: 0,
                    cc_load_policy: 1,
                    cc_lang_pref: 'id'
                },
                events: {
                    'onReady': function(){
                      if (textClicked==true){
                        onPlayerReady();
                      } else {
                        return;
                      };
                    },
                    'onStateChange': onPlayerStateChange
                }
            });
            },
        function() {
        console.log("ajax failure");
    }); 
    textClicked=true;
})});

function caption_detail(index){
    let subtitle_detail_area = document.getElementById("subtitle_detail_area");
    while (subtitle_detail_area.firstChild) {
        subtitle_detail_area.removeChild(subtitle_detail_area.firstChild);
    };
    let word =  content[index].textTokenized.word;
    let imi = content[index].textTokenized.imi;
    let s = '';
    for(let i = 0; i < word.length; i=i+1) {
        s = s + '<a href="/dic/' + word[i].slice(0,1) + '/' + word[i] + '"><p>' + word[i] + ' ' + imi[i] + '</p></a>';
    }
    subtitle_detail_area.innerHTML = s
};

function onPlayerReady() {
    player.loadVideoById({
      videoId: href,
      startSeconds: content[id].start_time *0.001
    })
};

function scrollSub(index){
    let topPos = subTitles.find('#' + index)["0"].parentNode.offsetTop;
    $(".panel").animate({
      scrollTop : topPos
    },250);
};

function colorSub(new_id) {
    if (s_id != -1) {
      $("#subtitle_text_" + s_id).css("background-color", "");
    }
    $("#subtitle_text_" + new_id).css("background-color", "#ccc");
    s_id = new_id;
};

function onPlayerStateChange(event) {
    clearTimeout(timer);
    if (textClicked==true && event.data==YT.PlayerState.PAUSED){
      textClicked=false;
      id=0;
      return;
    } else if (event.data==YT.PlayerState.PLAYING) {
      currentTime = player.getCurrentTime()* 1000;
      id = whichSub(currentTime);
      stopTime = content[id].end_time;
      if (subClicked){
        timeoutMillsec = stopTime - currentTime;
        if (repeat) {
          subClicked = false;
          timer = setTimeout(repeatVideo, timeoutMillsec);
        } else {
          subClicked = false;
          timer = setTimeout(pauseVideo, timeoutMillsec);
        }
      }else{
        timeoutMillsec = stopTime - currentTime;
        colorSub(id);
        scrollSub(id);
        caption_detail(id);
        if (textClicked==true) {
          timer = setTimeout(pauseVideo, timeoutMillsec);
          id=0;
        } else if (textClicked==false){
          timer = setTimeout(function() {sequenceSub(stopTime, id);}, timeoutMillsec);
        };
      };
    }
};

function whichSub(currentTime) {
    let index = 0;
    let time = content[index].end_time;
    while (time < currentTime) {
      index += 1;
      time = content[index].end_time;
    };
    return index
};
  
  // sequency change sub color by time
function sequenceSub(beforeendTime, index) {
    endTime = content[index].end_time;
    timeoutMillsec = endTime - beforeendTime;
    timer = setTimeout(function() {sequenceSub(endTime, index+1);}, timeoutMillsec);
    colorSub(index);
    scrollSub(index);
    caption_detail(index);
};

// pause video
function pauseVideo() {
    player.pauseVideo();
  };
  
  // repeat video
function repeatVideo() {
    player.pauseVideo();
    player.seekTo(startTime / 1000);
    player.playVideo();
};

$(document).on('click', '.subtitle', function () {
    startTime = content[this.id].start_time;
    subClicked = true;
    clearTimeout(timer);
    colorSub(this.id);
    scrollSub(this.id);
    caption_detail(this.id);
    player.pauseVideo();
    player.seekTo(startTime / 1000);
    player.playVideo();
  });