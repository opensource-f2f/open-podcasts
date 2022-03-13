waitForElementToDisplay(".widget-weixin",function(){
    removeByClassName('widget-weixin')
    removeByClassName('o_r_contact2')
    removeByClassName('a960_index')
    removeByClassName('a960_index a960_index_channel')
    document.getElementById('HMRichBox').remove()
},1000,19000);

function removeByClassName(name) {
    for (var i = 0; document.getElementsByClassName(name).length > 0;) {
        document.getElementsByClassName(name)[0].remove()
    }
}

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
