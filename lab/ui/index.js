const express = require('express');
const k8s = require("@kubernetes/client-node");
const {Client} = require("@kubernetes/client-node");
const crd = require("./crds/episodes.json"); //Import the express dependency
const app = express();              //Instantiate an express app, the main work horse of this server
const bodyParser = require('body-parser');
const YAML = require("yaml");
app.use(bodyParser());
const port = 5000;                  //Save the port number where your server will be listening

app.use(express.static('build'))
//Idiomatic expression in express to route and respond to a client request
app.get('/', (req, res) => {        //get requests to the root ("/") will route here
    res.sendFile('build/index.html', {root: __dirname});      //server responds by sending the index.html file to the client's browser
});

const https = require('follow-redirects').http
app.get('/stream/*', (req, res) => {
    const targetURL = req.url.replaceAll('/stream/', 'http://')
    https.get(targetURL, (rsp) => {
        for(let item in rsp.headers) {
            res.setHeader(item, rsp.headers[item])
        }
        rsp.pipe(res)
    })
})

app.get('/rsses', (req, res) => {
    const Client = require('kubernetes-client').Client
    const crd = require('./crds/rsses.json')
    const client = new Client({ version: '1.13' })
    client.addCustomResourceDefinition(crd)

    const all = client.apis['osf2f.my.domain'].v1alpha1.namespaces('default').rsses.get()
    all.then(response => {
        var space = 0
        if(req.query.pretty === 'true'){
            space = 2
        }
        res.end(JSON.stringify(response.body.items, undefined, space))
    })
})

app.get('/rsses/export', (req, res) => {
    const Client = require('kubernetes-client').Client
    const crd = require('./crds/rsses.json')
    const client = new Client({ version: '1.13' })
    const YAML = require('yaml');
    client.addCustomResourceDefinition(crd)

    const all = client.apis['osf2f.my.domain'].v1alpha1.namespaces('default').rsses.get()
    all.then(response => {
        const items = response.body.items
        const exportItems = []
        for(var i = 0; i < items.length; i++) {
            exportItems[i] = {
                apiVersion: items[i].apiVersion,
                kind: "RSS",
                metadata: {
                    name: items[i].metadata.name,
                },
                spec: {
                    address: items[i].spec.address,
                }
            }
        }

        res.set({"Content-Disposition":"attachment; filename=rsses.yaml"});

        var result = ""
        for(var i = 0; i < exportItems.length; i++) {
            const doc = new YAML.Document();
            doc.directivesEndMarker = true
            doc.contents = exportItems[i]
            result += doc.toString()
        }
        res.send(result);
    })
})

app.post('/rsses', (req, res) => {
    const Client = require('kubernetes-client').Client
    const crd = require('./crds/rsses.json')
    const client = new Client({ version: '1.13' })
    client.addCustomResourceDefinition(crd)

    client.apis['osf2f.my.domain'].v1alpha1.namespaces('default').rsses.post({
        body: {
            apiVersion: 'osf2f.my.domain/v1alpha1',
            kind: 'RSS',
            metadata: {
                generateName: (Math.random() + 1).toString(36).substring(7),
            },
            spec: {
                address: req.body.address
            }
        }
    })
    res.end('ok')
})

app.get('/episodes', (req, res) => {
    const Client = require('kubernetes-client').Client
    const crd = require('./crds/episodes.json')
    const crdRSS = require('./crds/rsses.json')
    const client = new Client({ version: '1.13' })
    client.addCustomResourceDefinition(crd)
    client.addCustomResourceDefinition(crdRSS)
    const rss = req.query.rss

    const rssObject = client.apis['osf2f.my.domain'].v1alpha1.namespaces('default').rsses(rss).get()
    rssObject.then(response => {
        const rssObj = response.body
        if (rssObj.status && rssObj.status.lastUpdateTime) {
            res.set({
                'Last-Modified' : rssObj.status.lastUpdateTime
            })
        }

        const all = client.apis['osf2f.my.domain'].v1alpha1.namespaces('default').episodes.get({ qs: { labelSelector: "rss=" + rss}})
        all.then(response => {
            var space = 0
            if(req.query.pretty === 'true'){
                space = 2
            }
            res.end(JSON.stringify(response.body.items, undefined, space))
        })
    })
});

app.get('/episodes/one', (req, res) => {
    const Client = require('kubernetes-client').Client
    const crd = require('./crds/episodes.json')
    const client = new Client({ version: '1.13' })
    client.addCustomResourceDefinition(crd)
    const name = req.query.name

    const all = client.apis['osf2f.my.domain'].v1alpha1.namespaces('default').episodes(name).get()
    all.then(response => {
        var space = 0
        if(req.query.pretty === 'true'){
            space = 2
        }
        res.end(JSON.stringify(response.body, undefined, space))
    })
})

app.get('/profiles', (req, res) => {
    const Client = require('kubernetes-client').Client
    const crd = require('./crds/profiles.json')
    const client = new Client({ version: '1.13' })
    client.addCustomResourceDefinition(crd)
    const name = req.query.name

    const all = client.apis['osf2f.my.domain'].v1alpha1.namespaces('default').profiles(name).get()
    all.then(response => {
        var space = 0
        if(req.query.pretty === 'true'){
            space = 2
        }
        res.end(JSON.stringify(response.body, undefined, space))
    })
});

app.post('/profiles/create', (req, res) => {
    const Client = require('kubernetes-client').Client
    const crd = require('./crds/profiles.json')
    const client = new Client({ version: '1.13' })
    client.addCustomResourceDefinition(crd)
    const name = req.query.name

    client.apis['osf2f.my.domain'].v1alpha1.namespaces('default').profiles.post({
        body: {
            apiVersion: 'osf2f.my.domain/v1alpha1',
            kind: 'Profile',
            metadata: {
                name: name,
            }
        }
    })
    res.end('ok')
})

app.delete('/profile/playLater', (req, res) => {
    const Client = require('kubernetes-client').Client
    const crd = require('./crds/profiles.json')
    const client = new Client({version: '1.13'})
    client.addCustomResourceDefinition(crd)

    profile = client.apis['osf2f.my.domain'].v1alpha1.namespaces('default').profiles(req.query.name).get()
    profile.then(response => {
        var found = false
        targetProfile = response.body
        if (!targetProfile.spec) {
            targetProfile.spec = {}
        }
        if (!targetProfile.spec.laterPlayList) {
            targetProfile.spec.laterPlayList = []
        }

        targetProfile.spec.laterPlayList.forEach(function (item, index){
            if(item.name === req.query.episode) {
                targetProfile.spec.laterPlayList.splice(index, 1)
                found = true
                return false
            }
        })
        if(found) {
            client.apis['osf2f.my.domain'].v1alpha1.namespaces('default').profiles(req.query.name)
                .put({
                    body:targetProfile,
                })
        }
    })
    res.status(200);
    res.end('ok')
})

app.post('/profile/playLater', (req, res) => {
    const Client = require('kubernetes-client').Client
    const crd = require('./crds/profiles.json')
    const client = new Client({ version: '1.13' })
    client.addCustomResourceDefinition(crd)

    profile = client.apis['osf2f.my.domain'].v1alpha1.namespaces('default').profiles(req.query.name).get()
    profile.then(response => {
        var found = false
        targetProfile = response.body
        if (!targetProfile.spec) {
            targetProfile.spec = {}
        }
        if (!targetProfile.spec.laterPlayList) {
            targetProfile.spec.laterPlayList = []
        }

        targetProfile.spec.laterPlayList.forEach(function (item, index){
            if(item.name === req.query.episode) {
                found = true
                return false
            }
        })
        if(!found) {
            targetProfile.spec.laterPlayList.push({
                name: req.query.episode
            })

            client.apis['osf2f.my.domain'].v1alpha1.namespaces('default').profiles(req.query.name)
                .put({
                    body:targetProfile,
                })
        }
    })
    res.status(200);
    res.end('ok')
})

app.post('/profile/playOver', (req, res) => {
    const Client = require('kubernetes-client').Client
    const crd = require('./crds/profiles.json')
    const client = new Client({ version: '1.13' })
    client.addCustomResourceDefinition(crd)

    const name = req.query.name
    const episode = req.query.episode

    profile = client.apis['osf2f.my.domain'].v1alpha1.namespaces('default').profiles(name).get()
    profile.then(response => {
        targetProfile = response.body
        targetProfile.spec.laterPlayList.forEach(function (item, index){
            if(item.name === episode) {
                targetProfile.spec.laterPlayList.splice(index, 1)
                if(targetProfile.spec.watchedList){
                    targetProfile.spec.watchedList.push({name: item.name})
                } else {
                    targetProfile.spec.watchedList = [{
                        name: item.name
                    }]
                }
                return false
            }
        })

        client.apis['osf2f.my.domain'].v1alpha1.namespaces('default').profiles(name)
            .put({
                body:targetProfile,
            })
    })
    res.status(200);
    res.end('ok')
})

app.post('/profile/social', (req, res) => {
    const Client = require('kubernetes-client').Client
    const crd = require('./crds/profiles.json')
    const client = new Client({ version: '1.13' })
    client.addCustomResourceDefinition(crd)

    const name = req.query.name
    const kind = req.query.kind
    const account = req.query.account

    profile = client.apis['osf2f.my.domain'].v1alpha1.namespaces('default').profiles(name).get()
    profile.then(response => {
        targetProfile = response.body
        if (!targetProfile.spec.socialLinks) {
            targetProfile.spec.socialLinks = {}
        }
        targetProfile.spec.socialLinks[kind] = account

        client.apis['osf2f.my.domain'].v1alpha1.namespaces('default').profiles(name)
            .put({
                body:targetProfile,
            })
    })
    res.status(200);
    res.end('ok')
})

app.listen(port, () => {            //server starts listening for any attempts from a client to connect at port: {port}
    console.log(`Now listening on port ${port}`);
});
