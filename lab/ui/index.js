const express = require('express');
const k8s = require("@kubernetes/client-node");
const {Client} = require("@kubernetes/client-node");
const crd = require("./crds/episodes.json"); //Import the express dependency
const app = express();              //Instantiate an express app, the main work horse of this server
const bodyParser = require('body-parser');
const YAML = require("yaml");
app.use(bodyParser());
const port = 5000;                  //Save the port number where your server will be listening

app.use(express.static('static'))
//Idiomatic expression in express to route and respond to a client request
app.get('/', (req, res) => {        //get requests to the root ("/") will route here
    res.sendFile('index.html', {root: __dirname});      //server responds by sending the index.html file to the client's browser
});

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
                name: req.body.name,
            },
            spec: {
                address: req.body.address
            }
        }
    })
    res.redirect('/')
})

app.get('/episodes', (req, res) => {
    const Client = require('kubernetes-client').Client
    const crd = require('./crds/episodes.json')
    const client = new Client({ version: '1.13' })
    client.addCustomResourceDefinition(crd)
    const rss = req.query.rss

    const all = client.apis['osf2f.my.domain'].v1alpha1.namespaces('default').episodes.get({ qs: { labelSelector: "rss=" + rss}})
    all.then(response => {
        var space = 0
        if(req.query.pretty === 'true'){
            space = 2
        }
        res.end(JSON.stringify(response.body.items, undefined, space))
    })
});

app.get('/profiles', (req, res) => {
    const Client = require('kubernetes-client').Client
    const crd = require('./crds/profiles.json')
    const client = new Client({ version: '1.13' })
    client.addCustomResourceDefinition(crd)

    const all = client.apis['osf2f.my.domain'].v1alpha1.namespaces('default').profiles.get()
    all.then(response => {
        var space = 0
        if(req.query.pretty === 'true'){
            space = 2
        }
        res.end(JSON.stringify(response.body.items, undefined, space))
    })
});

app.post('/profile/playLater', (req, res) => {
    const Client = require('kubernetes-client').Client
    const crd = require('./crds/profiles.json')
    const client = new Client({ version: '1.13' })
    client.addCustomResourceDefinition(crd)

    profile = client.apis['osf2f.my.domain'].v1alpha1.namespaces('default').profiles(req.query.name).get()
    profile.then(response => {
        var found = false
        response.body.spec.laterPlayList.forEach(function (item, index){
            if(item.name === req.query.episode) {
                found = true
                return false
            }
        })
        if(!found) {
            targetProfile = response.body
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

app.listen(port, () => {            //server starts listening for any attempts from a client to connect at port: {port}
    console.log(`Now listening on port ${port}`);
});
