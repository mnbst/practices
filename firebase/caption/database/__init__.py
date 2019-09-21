import os.path
from firebase_admin import credentials
from firebase_admin import firestore
import firebase_admin
import google.cloud.exceptions
import datetime
import urllib
import json
import isodate

dic = 'dictionary'
ng_dic = 'ng_dictionary'
video = 'video'
_caption = 'caption'
weat = 'usecase'

DEVELOPER_KEY = "AIzaSyArfj1GLTVUIyJ41bSm5I9gDLd5Frvn5Sk"

# connection to firestore
file_path = os.path.abspath(os.path.dirname(__file__))
path = os.path.join(file_path, 'mycredentialkey.json')
cred = credentials.Certificate(path)
firebase_admin.initialize_app(cred)
db = firestore.client()
batch = db.batch()


def dictionary_db(w, word_ini, meaning):
    doc = db.collection(dic).document(w)
    doc.set({
        u'word': w,
        u'word_ini': word_ini,
        u'word_imi': meaning,
    })


def dictionary_db_ng(w, word_ini):
    try:
        docng = db.collection(ng_dic).document(w)
        docng.set({
            u'word': w,
            u'word_ini': word_ini,
        })
        return
    except google.api_core.exceptions.InvalidArgument:
        return 


def get_word(w):
    return db.collection(dic).document(w).get()


def set_caption(href, search_result, title, script, element):
    batch.set(db.collection(video).document(href), {
        u'video_href': href,
        u'youtubeID': search_result["snippet"]["channelId"],
        u'video_img': search_result["snippet"]["thumbnails"]["medium"]["url"],
        u'video_title': title,
        u'video_update_time': search_result["snippet"]["publishedAt"],
        u'video_time': getduration(href),
        u'get_caption_ymd': datetime.datetime.today(),
    })

    batch.set(db.collection(_caption).document(href), {
        u'caption_id': href,
        u'caption_title': title,
    })

    for s in script:
        batch.set(db.collection(_caption).document(href).collection(u'content').document(), {
            u'index': s['index'],
            u'start_time': s['start_time'],
            u'end_time': s['end_time'],
            u'text': s['text'],
            u'text_tokenized': s['textTokenized'],
        })

    for t in element:
        batch.set(db.collection(weat).document(t), {
            href+u'_exists': True
        }, merge=True)
        for e in element[t]:
            batch.set(db.collection(weat).document(t).collection(u'cases').document(), {
                u'index': script[e]['index'],
                u'href': href,
                u'video_title': title,
                u'start_time': script[e]['start_time'],
                u'end_time': script[e]['end_time'],
                u'text': script[e]['text'],
                u'text_tokenized': script[e]['textTokenized'],
            })
        if len(batch._write_pbs) > 200:
            batch.commit()

    batch.commit()


def getduration(href):
    url = f"https://www.googleapis.com/youtube/v3/videos?id={href}&key={DEVELOPER_KEY}&part=contentDetails"
    response = urllib.request.urlopen(url).read()
    data = json.loads(response)
    duration = data['items'][0]['contentDetails']['duration']
    dur = isodate.parse_duration(duration).total_seconds()
    dur = int(dur)
    return '{:2}:{:02}'.format(dur//60, dur % 60)


def delete_href(notwant_and_delete):

    for nd in notwant_and_delete:
        for k in db.collection(_caption).document(nd).collection(u'content').stream():
            batch.delete(k.reference)
        k = None
        db.collection(video).document(nd).delete()
        db.collection(_caption).document(nd).delete()
        try:
            for k in db.collection(weat).where(nd+u'_exists', '==', True).stream():
                batch.update(k.reference, {
                    nd+u'_exists': firestore.DELETE_FIELD
                })
                for k2 in k.reference.collection(u'cases').where(u'href', '==', nd).stream():
                    batch.delete(k2.reference)
                    if len(batch._write_pbs) > 200:
                        batch.commit()
            k = None
            k2 = None
        except ValueError:
            continue

    batch.commit()

def video_data(href):
    return db.collection(video).document(href).get()

def caption_data(href):
    return db.collection(_caption).document(href).get()
