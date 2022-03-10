const express = require('express');
const k8s = require("@kubernetes/client-node");
const {Client} = require("@kubernetes/client-node");
const crd = require("./crds/episodes.json"); //Import the express dependency
const app = express();              //Instantiate an express app, the main work horse of this server
const bodyParser = require('body-parser');
const YAML = require("yaml");
const crdSub = require("./crds/subscriptions.json");
const toXML = require("to-xml").toXML
app.use(bodyParser());
const port = 5000;                  //Save the port number where your server will be listening

app.use(express.static('build'))
//Idiomatic expression in express to route and respond to a client request
app.get('/', (req, res) => {        //get requests to the root ("/") will route here
    res.sendFile('build/index.html', {root: __dirname});      //server responds by sending the index.html file to the client's browser
});

const https = require('follow-redirects').http
app.get('/stream/*', (req, res) => {
    try {
        const targetURL = req.url.replaceAll('/stream/', 'http://')
        https.get(targetURL, (rsp) => {
            for(let item in rsp.headers) {
                res.setHeader(item, rsp.headers[item])
            }
            rsp.pipe(res)
        })
    } catch (e) {
        console.log(e)
    }
})

const commandArgs = require('minimist')(process.argv.slice(2))
if (!commandArgs['defaultNamespace'] || commandArgs['defaultNamespace'] === "") {
    commandArgs['defaultNamespace'] = 'osf2f-system'
}
const defaultNamespace = commandArgs['defaultNamespace']

app.get('/namespaces', (req, res) => {
    const Client = require('kubernetes-client').Client
    const client = new Client({ version: '1.13' })

    client.api.v1.namespaces.get().then(resp => {
        res.end(JSON.stringify(resp.body.items, undefined, 2))
    })
})

app.get('/rsses', (req, res) => {
    const Client = require('kubernetes-client').Client
    const crd = require('./crds/rsses.json')
    const client = new Client({ version: '1.13' })
    client.addCustomResourceDefinition(crd)

    const all = client.apis['osf2f.my.domain'].v1alpha1.namespaces(defaultNamespace).rsses.get()
    all.then(response => {
        var space = 0
        if(req.query.pretty === 'true'){
            space = 2
        }
        res.end(JSON.stringify(response.body.items, undefined, space))
    })
})

app.get('/rsses/search', (req, res) => {
    const Client = require('kubernetes-client').Client
    const crd = require('./crds/rsses.json')
    const client = new Client({ version: '1.13' })
    client.addCustomResourceDefinition(crd)
    const category = req.query.category

    const all = client.apis['osf2f.my.domain'].v1alpha1.namespaces(defaultNamespace).rsses.get()
    all.then(response => {
        let items = response.body.items
        items = items.filter(function (val, index, arr) {
            if (val.spec.categories) {
                const result = val.spec.categories.filter(function (categoryName) {
                    return category === categoryName
                })
                return result.length > 0
            }
            return false
        })

        var space = 0
        if(req.query.pretty === 'true'){
            space = 2
        }
        res.end(JSON.stringify(items, undefined, space))
    })
})

app.get('/rsses/one', (req, res) => {
    const Client = require('kubernetes-client').Client
    const crd = require('./crds/rsses.json')
    const client = new Client({ version: '1.13' })
    client.addCustomResourceDefinition(crd)
    const name = req.query.name

    const all = client.apis['osf2f.my.domain'].v1alpha1.namespaces(defaultNamespace).rsses(name).get()
    all.then(response => {
        var space = 0
        if(req.query.pretty === 'true'){
            space = 2
        }
        res.end(JSON.stringify(response.body, undefined, space))
    })
})

app.get('/rsses/export', (req, res) => {
    const Client = require('kubernetes-client').Client
    const crd = require('./crds/rsses.json')
    const client = new Client({ version: '1.13' })
    const YAML = require('yaml');
    client.addCustomResourceDefinition(crd)

    const all = client.apis['osf2f.my.domain'].v1alpha1.namespaces(defaultNamespace).rsses.get()
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

app.get('/rsses/opml/export', (req, res) => {
    const Client = require('kubernetes-client').Client
    const crd = require('./crds/rsses.json')
    const client = new Client({ version: '1.13' })
    const YAML = require('yaml');
    client.addCustomResourceDefinition(crd)

    const all = client.apis['osf2f.my.domain'].v1alpha1.namespaces(defaultNamespace).rsses.get()
    all.then(response => {
        const items = response.body.items

        const bodyItems = []
        for(var i = 0; i < items.length; i++) {
            bodyItems.push({
                "@text": items[i].spec.title,
                "@title": items[i].spec.title,
                "@type": "rss",
                "@xmlUrl": items[i].spec.address,
                "@htmlUrl": ""
            })
        }
        const opml = {
            "?xml version='1.0' encoding='UTF-8' standalone='no' ?": null,
            "opml": {
                "@version": "2.0",
                head: {
                    title: 'From Open Podcasts',
                    dateCreated: new Date().toISOString(),
                },
                body: {
                    outline: bodyItems
                }
            }
        }

        res.set({"Content-Disposition":"attachment; filename=rsses.opml"});

        res.send(toXML(opml, null, 2));
    })
})

app.post('/rsses', (req, res) => {
    const Client = require('kubernetes-client').Client
    const crd = require('./crds/rsses.json')
    const client = new Client({ version: '1.13' })
    client.addCustomResourceDefinition(crd)

    client.apis['osf2f.my.domain'].v1alpha1.namespaces(defaultNamespace).rsses.post({
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

    const rssObject = client.apis['osf2f.my.domain'].v1alpha1.namespaces(defaultNamespace).rsses(rss).get()
    rssObject.then(response => {
        const rssObj = response.body
        if (rssObj.status && rssObj.status.lastUpdateTime) {
            res.set({
                'Last-Modified' : rssObj.status.lastUpdateTime
            })
        }

        const all = client.apis['osf2f.my.domain'].v1alpha1.namespaces(defaultNamespace).episodes.get({ qs: { labelSelector: "rss=" + rss}})
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

    const all = client.apis['osf2f.my.domain'].v1alpha1.namespaces(defaultNamespace).episodes(name).get()
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

    var space = 0
    if(req.query.pretty === 'true'){
        space = 2
    }

    try {
        const all = client.apis['osf2f.my.domain'].v1alpha1.namespaces(defaultNamespace).profiles(name).get()
        all.then(response => {
            res.end(JSON.stringify(response.body, undefined, space))
        }).catch(err => {
            console.error('Error: ', err)
            res.send(JSON.stringify('{message:"error"}', undefined, space))
        })
    } catch (e) {
        console.error('Error: ', e)
        res.send(JSON.stringify('{message:"error"}', undefined, space))
    }
});

app.post('/profiles/create', (req, res) => {
    const Client = require('kubernetes-client').Client
    const crd = require('./crds/profiles.json')
    const client = new Client({ version: '1.13' })
    client.addCustomResourceDefinition(crd)
    const name = req.query.name

    try {
        client.apis['osf2f.my.domain'].v1alpha1.namespaces(defaultNamespace).profiles.post({
            body: {
                apiVersion: 'osf2f.my.domain/v1alpha1',
                kind: 'Profile',
                metadata: {
                    name: name,
                }
            }
        }).error(res => {
            console.log(res)
        })
        res.end('ok')
    } catch (e) {
        console.error('Error: ', e)
    }
})

app.delete('/profile/playLater', (req, res) => {
    const Client = require('kubernetes-client').Client
    const crd = require('./crds/profiles.json')
    const client = new Client({version: '1.13'})
    client.addCustomResourceDefinition(crd)

    profile = client.apis['osf2f.my.domain'].v1alpha1.namespaces(defaultNamespace).profiles(req.query.name).get()
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
            client.apis['osf2f.my.domain'].v1alpha1.namespaces(defaultNamespace).profiles(req.query.name)
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

    profile = client.apis['osf2f.my.domain'].v1alpha1.namespaces(defaultNamespace).profiles(req.query.name).get()
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

            client.apis['osf2f.my.domain'].v1alpha1.namespaces(defaultNamespace).profiles(req.query.name)
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

    profile = client.apis['osf2f.my.domain'].v1alpha1.namespaces(defaultNamespace).profiles(name).get()
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

        client.apis['osf2f.my.domain'].v1alpha1.namespaces(defaultNamespace).profiles(name)
            .put({
                body:targetProfile,
            })
    })
    res.status(200);
    res.end('ok')
})

app.post('/profile/notifier', (req, res) => {
    const Client = require('kubernetes-client').Client
    const crd = require('./crds/profiles.json')
    const crdNotifiers = require('./crds/notifiers.json')
    const client = new Client({ version: '1.13' })
    client.addCustomResourceDefinition(crd)
    client.addCustomResourceDefinition(crdNotifiers)

    const name = req.query.name
    const url = req.query.url

    profile = client.apis['osf2f.my.domain'].v1alpha1.namespaces(defaultNamespace).profiles(name).get()
    profile.then(response => {
        targetProfile = response.body
        if (!targetProfile.spec.notifier) {
            client.apis['osf2f.my.domain'].v1alpha1.namespaces(defaultNamespace).notifiers.post({
                body: {
                    apiVersion: 'osf2f.my.domain/v1alpha1',
                    kind: 'Notifier',
                    metadata: {
                        generateName: name,
                    },
                    spec: {
                        feishu: {
                            webhook_url: url,
                        },
                    }
                }
            }).then(res => {
                targetProfile.spec.notifier = {
                    name: res.body.metadata.name,
                }
                client.apis['osf2f.my.domain'].v1alpha1.namespaces(defaultNamespace).profiles(name)
                    .put({
                        body: targetProfile,
                    })
            })
        } else {
            const notifierName = targetProfile.spec.notifier.name
            const notifier = client.apis['osf2f.my.domain'].v1alpha1.namespaces(defaultNamespace).notifiers(notifierName).get()
            notifier.then(notifierRes => {
                const notifierObj = notifierRes.body
                notifierObj.spec.feishu.webhook_url = url

                client.apis['osf2f.my.domain'].v1alpha1.namespaces(defaultNamespace).notifiers(notifierName)
                    .put({
                        body: notifierObj,
                    })
            })
        }
    })
    res.status(200);
    res.end('ok')
})

app.get('/notifiers/one', (req, res) => {
    const Client = require('kubernetes-client').Client
    const crd = require('./crds/notifiers.json')
    const client = new Client({ version: '1.13' })
    client.addCustomResourceDefinition(crd)
    const name = req.query.name

    const all = client.apis['osf2f.my.domain'].v1alpha1.namespaces(defaultNamespace).notifiers(name).get()
    all.then(response => {
        var space = 0
        if(req.query.pretty === 'true'){
            space = 2
        }
        res.end(JSON.stringify(response.body, undefined, space))
    })
})

app.post('/profile/social', (req, res) => {
    const Client = require('kubernetes-client').Client
    const crd = require('./crds/profiles.json')
    const client = new Client({ version: '1.13' })
    client.addCustomResourceDefinition(crd)

    const name = req.query.name
    const kind = req.query.kind
    const account = req.query.account

    profile = client.apis['osf2f.my.domain'].v1alpha1.namespaces(defaultNamespace).profiles(name).get()
    profile.then(response => {
        targetProfile = response.body
        if (!targetProfile.spec.socialLinks) {
            targetProfile.spec.socialLinks = {}
        }
        targetProfile.spec.socialLinks[kind] = account

        client.apis['osf2f.my.domain'].v1alpha1.namespaces(defaultNamespace).profiles(name)
            .put({
                body:targetProfile,
            })
    })
    res.status(200);
    res.end('ok')
})

app.post('/subscribe', (req, res) => {
    const Client = require('kubernetes-client').Client
    const crd = require('./crds/profiles.json')
    const crdSub = require('./crds/subscriptions.json')
    const client = new Client({ version: '1.13' })
    client.addCustomResourceDefinition(crd)
    client.addCustomResourceDefinition(crdSub)

    const profileName = req.query.profile
    const rss = req.query.rss

    profile = client.apis['osf2f.my.domain'].v1alpha1.namespaces(defaultNamespace).profiles(profileName).get()
    profile.then(response => {
        targetProfile = response.body
        if (!targetProfile.spec.subscription) {
            client.apis['osf2f.my.domain'].v1alpha1.namespaces(defaultNamespace).subscriptions.post({
                body: {
                    apiVersion: 'osf2f.my.domain/v1alpha1',
                    kind: 'Subscription',
                    metadata: {
                        generateName: profileName,
                    },
                    spec: {
                        rssList: [{
                            name: rss,
                        }],
                    }
                }
            }).then(res => {
                targetProfile.spec.subscription = {
                    name: res.body.metadata.name,
                }
                client.apis['osf2f.my.domain'].v1alpha1.namespaces(defaultNamespace).profiles(profileName)
                    .put({
                        body: targetProfile,
                    })
            })
        } else {
            const subName = targetProfile.spec.subscription.name
            const sub = client.apis['osf2f.my.domain'].v1alpha1.namespaces(defaultNamespace).subscriptions(subName).get()
            sub.then(subRes => {
                subscription = subRes.body
                if (subscription.spec.rssList) {
                    subscription.spec.rssList.push({
                        name: rss,
                    })
                } else {
                    subscription.spec.rssList = [{
                        name: rss,
                    }]
                }

                client.apis['osf2f.my.domain'].v1alpha1.namespaces(defaultNamespace).subscriptions(subName)
                    .put({
                        body:subscription,
                    })
            })
        }
    })
    res.status(200);
    res.end('ok')
})

app.post('/unsubscribe', (req, res) => {
    const Client = require('kubernetes-client').Client
    const crd = require('./crds/profiles.json')
    const crdSub = require('./crds/subscriptions.json')
    const client = new Client({ version: '1.13' })
    client.addCustomResourceDefinition(crd)
    client.addCustomResourceDefinition(crdSub)

    const profileName = req.query.profile
    const rss = req.query.rss

    profile = client.apis['osf2f.my.domain'].v1alpha1.namespaces(defaultNamespace).profiles(profileName).get()
    profile.then(response => {
        targetProfile = response.body
        if (targetProfile.spec.subscription) {
            const subName = targetProfile.spec.subscription.name
            const sub = client.apis['osf2f.my.domain'].v1alpha1.namespaces(defaultNamespace).subscriptions(subName).get()
            sub.then(subRes => {
                subscription = subRes.body
                if (subscription.spec.rssList) {
                    let found = false
                    subscription.spec.rssList.forEach(function (item, index){
                        if(item.name === rss) {
                            subscription.spec.rssList.splice(index, 1)
                            found = true
                            return false
                        }
                    })

                    if (found) {
                        client.apis['osf2f.my.domain'].v1alpha1.namespaces(defaultNamespace).subscriptions(subName)
                            .put({
                                body:subscription,
                            })
                    }
                }
            })
        }
    })
    res.status(200);
    res.end('ok')
})

app.get('/subscriptions/one', (req, res) => {
    const Client = require('kubernetes-client').Client
    const crd = require('./crds/subscriptions.json')
    const client = new Client({ version: '1.13' })
    client.addCustomResourceDefinition(crd)
    const name = req.query.name

    const all = client.apis['osf2f.my.domain'].v1alpha1.namespaces(defaultNamespace).subscriptions(name).get()
    all.then(response => {
        var space = 0
        if(req.query.pretty === 'true'){
            space = 2
        }
        res.end(JSON.stringify(response.body, undefined, space))
    })
})

app.listen(port, () => {            //server starts listening for any attempts from a client to connect at port: {port}
    console.log(`Now listening on port ${port}`);
});
