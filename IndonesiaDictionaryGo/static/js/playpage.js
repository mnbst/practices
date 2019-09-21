// some var
var player;
var timeoutMillsec = 0;
var timer;
var stopTime = 0;
var startTime = 0;
var subClicked = false;
var v_id;
var s_id = -1; // colored sub id
var url;
var repeat = false;
var start = 0;

// get video href
url = window.location.href;
url = url.replace('#', '');
for (var i = url.length, ref = url.length; i >= 0; i = i - 1) {
  if (url[i] == '/') {
    v_id = url.substring(i + 1, ref);
    break
  }
};

// load video's subtitles
var obj =document.getElementById("json-data");
var json=obj.textContent;
content = JSON.parse(json);

$(window).on('load', function () {
  content.forEach(function (element, index) {
    $("#subtitle_block").append("<tr><td>" + "<a id='" + index + "'' class='subtitle' href='#' >" + "<i class='far fa-play-circle fa-lg'></i>" + "</a></td>" + "<td><a id='" + index + "'' class='subtitle' href='#' ><span id='subtitle_text_" + index + "' >" + element["text"] + "</span></a></td></tr>");
  });
});

// initial youtube iframe api
var tag = document.createElement('script');
tag.src = "http://www.youtube.com/player_api";
var firstScriptTag = document.getElementsByTagName('script')[0];
firstScriptTag.parentNode.insertBefore(tag, firstScriptTag);


function onYouTubePlayerAPIReady() {
  player = new YT.Player('player', {
    height: '390',
    width: '640',
    videoId: v_id,
    playerVars: {
      autoplay: 0,
      modestbranding: 0,
      rel: 0,
      cc_load_policy: 1,
      cc_lang_pref: 'id'
    },
    events: {
      //'onReady': onPlayerReady,
      'onStateChange': onPlayerStateChange
    }
  });
};

function onPlayerStateChange(event) {
  clearTimeout(timer);
  if (event.data == YT.PlayerState.PLAYING) {
    var currentTime = player.getCurrentTime()* 1000;
    var id = whichSub(currentTime);
    stopTime = content[id]["end_time"];
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
      timer = setTimeout(function() {sequenceSub(stopTime, id);}, timeoutMillsec);
      colorSub(id);
      scrollSub(id);
      caption_detail(id);
    };
  }
};

// サブタイトルクリックしたときの処理
$(document).on('click', '.subtitle', function () {
  startTime = content[this.id]["start_time"]
  stopTime = content[this.id]["end_time"]
  subClicked = true;
  colorSub(this.id);
  scrollSub(this.id);
  caption_detail(this.id);
  clearTimeout(timer);
  player.pauseVideo();
  player.seekTo(startTime / 1000);
  player.playVideo();
});

function caption_detail(id){
  // var subtitle_area = document.getElementByid("subtitle_area");
  var textTokenized = content[id]["textTokenized"];
  var subtitle_detail_area = document.getElementById("subtitle_detail_area");
  while (subtitle_detail_area.firstChild) {
    subtitle_detail_area.removeChild(subtitle_detail_area.firstChild);
  }
  var word =  textTokenized["word"];
  var imi = textTokenized["imi"];
  var s = '';
  for(let i = 0; i < word.length; i++) {
    s = s + '<a href="/dic/' + word[i].slice(0,1) + '/' + word[i] + '"><li>' + word[i] + ' ' + imi[i] + '</li></a>';
  }
  subtitle_detail_area.innerHTML = s
};

// handle subtitle's color
function colorSub(new_id) {
  if (s_id != -1) {
    $("#subtitle_text_" + s_id).css("background-color", "");
  }
  $("#subtitle_text_" + new_id).css("background-color", "#ccc");
  s_id = new_id;
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

// find sub id by time 
function whichSub(currentTime) {
  var id = 0;
  var time = content[id]["end_time"];
  while (time < currentTime) {
    id += 1;
    time = content[id]["end_time"];
  };
  return id
};

// sequency change sub color by time
function sequenceSub(beforeendTime, id) {
  //endTime = content[id-1]["end-time"];
  endTime = content[id]["end_time"];
  timeoutMillsec = endTime - beforeendTime;
  timer = setTimeout(function() {sequenceSub(endTime, id+1);}, timeoutMillsec);
  colorSub(id);
  scrollSub(id);
  caption_detail(id);
};

// scroll sub to top
function scrollSub(id){
  var topPos = document.getElementById(id).parentNode.parentNode.offsetTop;
  $(".panel").animate({
    scrollTop : topPos
  },250)
};
