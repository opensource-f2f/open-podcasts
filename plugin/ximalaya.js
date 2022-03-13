waitForElementToDisplay(".controls",function(){
    let eles = document.getElementsByClassName('controls')
    if (eles.length > 0) {
        let rss = document.createElement('div')
        rss.innerHTML = 'RSS: <a target="_blank" href="'  + window.location.toString() + '.xml">' + window.location.toString() + '.xml</a>'
        eles[0].parentElement.append(rss)
    }
},1000,19000);

function waitForElementToDisplay(selector, callback, checkFrequencyInMs, timeoutInMs) {
    var startTimeInMs = Date.now();
    (function loopSearch() {
        if (document.querySelector(selector) != null) {
            callback();
            return;
        }
        else {
            setTimeout(function () {
                if (timeoutInMs && Date.now() - startTimeInMs > timeoutInMs)
                    return;
                loopSearch();
            }, checkFrequencyInMs);
        }
    })();
}
