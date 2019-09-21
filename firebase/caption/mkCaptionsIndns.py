# -*- coding: utf-8 -*-
from googleapiclient.discovery import build
from googleapiclient.errors import HttpError
from oauth2client.tools import argparser
import requests
import xml.etree.ElementTree as ET
import nltk
import re
import google.cloud.exceptions
from html import unescape
from google.cloud.firestore_v1 import ArrayUnion, ArrayRemove
from bs4 import BeautifulSoup
import urllib
import database

# Set DEVELOPER_KEY to we API key value from worde APIs & auword > Registered apps
# tab of
#   https://cloud.google.com/console
# ＡＰＩデベロッパーキーを入れる
DEVELOPER_KEY = "AIzaSyArfj1GLTVUIyJ41bSm5I9gDLd5Frvn5Sk"
YOUTUBE_API_SERVICE_NAME = "youtube"
YOUTUBE_API_VERSION = "v3"

# クロールするページ数
page_to_crowl = 10
# １動画ごとの表示言語数
count_of_languagecode = 3
# １動画ごとの最小文数。空欄のやつとかはじける。
minimum_count_of_sentences = 10
# １ページあたりの動画取得数。上限50
video_per_page = 50
# 除外したいyoutubeID。実行時にＤＢにあるやつも消す
delete = []
renewal_word = []
notwant = ['-3fl9HVCk3o', 'vN_5MBIxjcs', 'KLFjdzvTLWM', 'CS9am2ctk1M', 'r9QG3eFJUJk', '1jQdXHdJIQw', 'FCeZqjTwO2M',
           '_lhS0oilP1Q', 'auNX2WZLuXA', 'AFD0nFg2s4M', 'jvfAEtDlR2A', 'C-vQZU1xAf4', 'U8Fd0hqLApU', 'OvMN_8ClEgA', 'R2JhxW1rT8I',
           'pRpNiYI0g64', 'xqJ-ag2VLbk', 'SDD0r15AprU', 'rEGmWKbHUAY', 'UlFS0utBxB0', 'APTRWJ8Q1UU', 'lAeAO51aVKA', 'Pl-RYkIXPnk',
           'e2RBRUFwmuw', 'lIUYItKBbBc', 'SJTzMHDaOlg', 'fX6sg2xuvCw', 'JT3Wc5alyfQ', 'FlBdDhexQRc', 'GePukn4vRV8', '-Xk5aUvZCc0',
           '5H-41WPSsHs', 'bQ_WHZ_WS3s', '5M1_ufFxvlE', '1l_Y8X9HhUA', 'd3yoDVp0hzc', 'I7bVw0nGarQ', 'XPzB_GSMABo', '_7SNK6m8okE',
           'sclwgAXiHhI', 'jh8dEnVTvs8', 'Xw5RBrqCoL8', 'm37ByQITb9k', 'BMmnrdlxdbg', 'DCSC4OG4nEE', 'n6hJXjka4Xk', 'OCA9alkiptw',
           'ah0Vcf1X8r0', 'EPgOhNeZYOc', 'qGpvtvSzGQ8', 'si16hSJF6v8', 'A3LRdMUzI-4', 'KaXx2z5rYIA', 'MMIABcbEIGg', '9usVvl4q2Dc',
           '11I-WD-yBDU']

database.delete_href(delete + notwant)

# 除外する記号
regex = r"[!-,.-@[-`{-~]"


def youtube_search(options):
    n = 0

    youtube = build(YOUTUBE_API_SERVICE_NAME,
                    YOUTUBE_API_VERSION, developerKey=DEVELOPER_KEY)

    # Call worde search.list mewordod to retrieve results matching worde specified
    search_response = youtube.search().list(
        q=options.q,
        part="id,snippet",
        maxResults=options.max_results,
        regionCode="id",
        relevanceLanguage="id",
        safeSearch="moderate",
        type="video",
        videoCaption="closedCaption",
        videoSyndicated="true"
    ).execute()

    pagetoken = search_response.get("nextPageToken", str)
    n = fillindb(search_response)
    print(n)

    for _ in range(0, page_to_crowl-1):
        search_response.update(
            youtube.search().list(
                q=options.q,
                pageToken=pagetoken,
                part="id,snippet",
                maxResults=options.max_results,
                regionCode="id",
                relevanceLanguage="id",
                safeSearch="moderate",
                type="video",
                videoCaption="closedCaption",
                videoSyndicated="true"
            ).execute()
        )
        pagetoken = search_response.get("nextPageToken", str)
        n = fillindb(search_response)

        print(n)


def fillindb(response):

    n = 0
    # Add each result to worde appropriate list, and worden display worde lists of
    # matching videos, channels, and playlists.
    for search_result in response.get("items", []):
        if search_result["id"]["kind"] == "youtube#video":
            href = search_result["id"]["videoId"]
            if href in notwant:
                continue
        else:
            href = ''
            continue

        vdata = None
        cdata = None
        try:
            vdata = database.video_data(href)
        except google.cloud.exceptions.NotFound:
            vdata is None
        try:
            cdata = database.caption_data(href)
        except google.cloud.exceptions.NotFound:
            cdata is None

        if vdata._exists == True and cdata._exists == True:
            continue

        script = []
        langcodes = []
        url = f'https://www.youtube.com/api/timedtext?type=list&v={href}'
        req = requests.get(url)
        try:
            root = ET.fromstring(req.content)
        except ET.ParseError:
            continue

         # 字幕言語取得
        for r in root.iter():
            lancd = r.attrib.get('lang_code')
            langcodes.append(lancd)

        # 字幕言語の数を指定してる
        if len(langcodes) < count_of_languagecode+1 and 'id' in langcodes:
            name = r.attrib.get('name')
            script, element = mkscript(name, href)

            if len(script) < minimum_count_of_sentences:
                continue

            n += 1
            title = search_result["snippet"]["title"]
            title = unescape(title)
            #caption=json.dumps(script, ensure_ascii=False)

            #set to firestore
            database.set_caption(href, search_result, title, script, element)
            print(href)

    return n


def mkscript(name, href):
    url = f"http://video.google.com/timedtext?lang=id&name={name}&v={href}"
    req = requests.get(url)
    root = ET.fromstring(req.content)

    # 字幕つくる。構成単語も
    element = {}
    script = []
    index = 0
    for r in root.iter():

        word = []
        imi = []
        try:
            text = re.sub("\n", " ", r.text)
            text = unescape(text)
            print(text)
        except TypeError:
            text = r.text
        except AttributeError:
            text = r.text

        try:
            row = nltk.tokenize.word_tokenize(text)
            row = [re.sub(regex, '', item) for item in row]
            word = [i for i in row if not i == '']
        except TypeError:
            word = []
            continue

        if not word:
            continue

        start = r.attrib.get('start')
        try:
            start_time = int(float(start)*1000)
        except TypeError:
            start_time = 0

        dur = r.attrib.get('dur')
        try:
            end_time = start_time+int(float(dur)*1000)
        except TypeError:
            end_time = start_time

        # DBから意味を取得。なければinsert
        for j, w in enumerate(word):
            meaning = ''
            w = w.lower()
            if w.startswith('-') or w.endswith('-'):
                w = re.sub('-', '', w)

            meaning = getimi(w)

            if meaning != '' and w in element:
                element[w].append(index)
            elif meaning != '' and w not in element:
                element[w] = [index]

            if meaning == '':
                word[j] = ''
            else:
                imi.append(meaning)
                word[j] = w

        word = [i for i in word if not i == '']
        token = {"word": [], "imi": []}
        token["word"] = word
        token["imi"] = imi
        textTokenized = token

        sentence = {"index": 0, "href": href, "start_time": 0,
                    "end_time": 0, "text": '', "textTokenized": textTokenized}
        sentence["index"] = index
        sentence["start_time"] = start_time
        sentence["end_time"] = end_time
        sentence["text"] = text

        script.append(sentence)
        index += 1

    return script, element


def getimi(w):
    meaning = ''
    meaning = getimifromdb(w)
    if meaning != '':
        return meaning

    url = f'https://njjn.weblio.jp/content/{w}'
    r = requests.get(url)
    soup = BeautifulSoup(r.text, 'lxml')

    try:
        elmnts_crslnk = soup.find_all(class_="crosslink")
        elmnts_igngj = soup.find_all(class_="Igngj")
        if elmnts_crslnk != [] or elmnts_igngj != []:
            meaning = gettext(elmnts_crslnk+elmnts_igngj)
    except urllib.error.HTTPError:
        meaning = ""
    except urllib.error.URLError:
        meaning = ""

    # 接頭語・接尾語とか　　
    if not meaning and w.endswith('nya'):
        w1 = w[:len(w)-3]
        meaning = getimi(w1)
        if meaning != '':
            meaning = meaning+"(nya=特定の事柄・人を表す接尾辞)"
            return meaning
    elif not meaning and w.endswith('kah'):
        w1 = w[:len(w)-3]
        meaning = getimi(w1)
        if meaning != '':
            meaning = meaning+"(kah=～ですか？)"
            return meaning
    elif not meaning and w.endswith('an') and not w.endswith('kan'):
        w1 = w[:len(w)-2]
        meaning = getimi(w1)
        if meaning != '':
            meaning = meaning+"(an=単位/内容を特定する接尾辞)"
    elif not meaning and w.endswith('in') and not '-' in w:
        w1 = w[:len(w)-2]
        meaning = getimi(w1)
        if meaning != '':
            meaning = meaning+"(in=ジャカルタ方言・他動詞の語幹をつくる接尾辞)"
    elif not meaning and w.endswith('i'):
        w1 = w[:len(w)-1]
        meaning = getimi(w1)
        if meaning != '':
            meaning = meaning+"(i=前置詞を代替する接尾辞/動作の反復・集中)"
    elif not meaning and w.endswith('kan'):
        w1 = w[:len(w)-3]
        meaning = getimi(w1)
        if meaning != '':
            meaning = meaning+"(kan=他動詞の語幹をつくる接尾辞)"
    elif not meaning and w.startswith('mu'):
        w1 = w[2:]
        meaning = getimi(w1)
        if meaning != '' and not '=' in meaning:
            meaning = meaning+"(mu=kamu)"
            return meaning
    elif not meaning and w.startswith('ku'):
        w1 = w[2:]
        meaning = getimi(w1)
        if meaning != '' and not '=' in meaning:
            meaning = meaning+"(ku=aku)"
            return meaning
    elif not meaning and w.endswith('lah'):
        w1 = w[:len(w)-3]
        meaning = getimi(w1)
        if meaning != '':
            meaning = meaning+"(lah=強調)"
            return meaning
    elif not meaning and w.startswith('di'):
        w1 = w[2:]
        meaning = getimi(w1)

        if meaning == '':
            if w1.startswith(('l', 'r', 'm', 'n', 'w', 'y')):
                w2 = 'me'+w1
                meaning = getimi(w2)
            elif w1.startswith('t'):
                w2 = 'men'+w1[1:]
                meaning = getimi(w2)
            elif w1.startswith('p'):
                w2 = 'mem'+w1[1:]
                meaning = getimi(w2)
            elif w1.startswith(('c', 'j', 'z', 'sy', 'd')):
                w2 = 'men'+w1
                meaning = getimi(w2)
            elif w1.startswith('b'):
                w2 = 'mem'+w1
                meaning = getimi(w2)
            elif w1.startswith('s'):
                w2 = 'meny'+w1[1:]
                meaning = getimi(w2)
            else:
                w2 = 'meng'+w1
                meaning = getimi(w2)

        if meaning != '':
            meaning = meaning+"(di=受け身を表す接頭辞/～で、～に)"
            return meaning
    elif not meaning and w.endswith('pun'):
        w1 = w[:len(w)-3]
        meaning = getimi(w1)
        if meaning != '':
            meaning = meaning+"(pun=～でさえ/～でも)"
            return meaning
    elif not meaning and w.startswith('ber'):
        w1 = w[3:]
        meaning = getimi(w1)
        if meaning != '':
            meaning = meaning+"(ber=～を持っている／身につけている／伴っている)"
            return meaning
    elif not meaning and w.endswith('mu'):
        w1 = w[:len(w)-2]
        meaning = getimi(w1)
        if meaning != '':
            meaning = meaning+"(mu=kamu)"
            return meaning
    elif not meaning and w.endswith('ku'):
        w1 = w[:len(w)-2]
        meaning = getimi(w1)
        if meaning != '':
            meaning = meaning+"(ku=aku)"
            return meaning
    elif not meaning and w.startswith('ter'):
        w1 = w[3:]
        meaning = getimi(w1)
        if meaning != '':
            meaning = meaning+"(ter=最も～/動作対象＋～してしまう，～してしまっている)"
    elif not meaning and w.startswith('se'):
        w1 = w[2:]
        meaning = getimi(w1)
        if meaning != '':
            meaning = meaning+"(se=1つの/1回の/同じ/全体)"
    elif not meaning and '-' in w:
        w2 = w.split('-')[1]
        w1 = w.split('-')[0]
        if w2 != '' and w2 in w1:
            meaning = getimi(w1)
        if meaning != '':
            meaning = meaning+"(\"-\"=複数/動作の繰り返し)"
            return meaning

    if w.startswith('-') or w.endswith('-'):
        w = re.sub('-', '', w)
    word_ini = w[0:1]

    if meaning != '':
        meaning = meaning.strip()
        #set word
        database.dictionary_db(w, word_ini, meaning)
    elif meaning == '':
        #set ng word
        database.dictionary_db_ng(w, word_ini)   

    return meaning         

def getimifromdb(w):
    try:
        doc = database.get_word(w)
    except google.api_core.exceptions.InvalidArgument:
        meaning = ""
        return meaning

    if doc._exists == True:
        meaning = doc._data['word_imi']
        return meaning
    else:
        meaning = ''
    return meaning


def gettext(elmnts):
    contents = []
    for tag in elmnts:
        contents = contents + tag.text.strip().split(',')
    c_unique = list(set(contents))
    return ','.join(c_unique)


if __name__ == "__main__":
    argparser.add_argument("--q", help="Search term", default="indonesia")
    argparser.add_argument(
        "--max-results", help="Max results", default=video_per_page)
    args = argparser.parse_args()

    try:
        youtube_search(args)
    except HttpError as e:
        print("An HTTP error %d occurred:\n%s" % (e.resp.status, e.content))
