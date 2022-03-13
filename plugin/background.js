chrome.contextMenus.create({
    "title" : "Collect RSS",
    "onclick" : collect
});

function collect(){
    if (window.location.toString().startsWith('https://www.ximalaya.com/album/')) {
    }
}
