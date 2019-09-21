import * as functions from 'firebase-functions';
import * as admin from 'firebase-admin';
admin.initializeApp({ credential: admin.credential.applicationDefault() });

export const thumbnails =
    functions.region('asia-northeast1').https.onRequest(async (_, response) => {
        response.header('Content-Type', 'application/json')
        try {
            const snapshot = await admin.firestore().collection("video").get();
            response.send(snapshot.docs.map(doc => doc.data()));
        }
        catch (error) {
            console.log(error)
            response.status(500).send(error)
        }
    });

export const caption =
    functions.region('asia-northeast1').https.onRequest(async (request, response) => {
        response.header('Content-Type', 'application/json')
        try {
            const snapshot = await admin.firestore().doc('caption/' + request.query.id).collection('content').orderBy('index', 'asc').get();
            response.send(snapshot.docs.map(doc => doc.data()));
        }
        catch (error) {
            console.log(error)
            response.status(500).send(error)
        }
    });

export const meaning =
    functions.region('asia-northeast1').https.onRequest(async (request, response) => {
        response.header('Content-Type', 'application/json')
        try {
            const snapshot = await admin.firestore().doc('dictionary/' + request.query.word).get()
            response.send(snapshot.data());
        }
        catch (error) {
            console.log(error)
            response.status(500).send(error)
        }
    });

export const dictionary =
    functions.region('asia-northeast1').https.onRequest(async (request, response) => {
        response.header('Content-Type', 'application/json')
        try {
            const snapshot = await admin.firestore().collection('dictionary').where('word_ini', '==', request.query.initial).get()
            response.send(snapshot.docs.map(doc => doc.data()));
        }
        catch (error) {
            console.log(error)
            response.status(500).send(error)
        }
    });

export const usecase =
    functions.region('asia-northeast1').https.onRequest(async (request, response) => {
        response.header('Content-Type', 'application/json')
        try {
            const snapshot = await admin.firestore().doc('usecase/' + request.query.word).collection('cases').get()
            response.send(snapshot.docs.map(doc => doc.data()));
        }
        catch (error) {
            console.log(error)
            response.status(500).send(error)
        }
    });


