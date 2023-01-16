function GetBrowserData() {
    const browserData = {}
    if (window) {
        if (window.screen) {
            browserData.browserColorDepth = window.screen.colorDepth;
            browserData.browserScreenHeight = window.screen.height;
            browserData.browserScreenWidth = window.screen.width;
        }
        if (window.navigator) {
            browserData.browserUserAgent = window.navigator.userAgent;
            browserData.browserJavaEnabled = window.navigator.javaEnabled();
            browserData.browserLanguage = window.navigator.language || window.navigator.browserLanguage || window.navigator.userLanguage;
        }
    }
    browserData.browserTZ = (new Date()).getTimezoneOffset();
    browserData.browserJavascriptEnabled = true;
    return browserData;
}

function SendMethodRequest(threeDSMethodURL, threeDSServerTransID, methodNotificationURL) {
    let frameContainer = document.getElementById('methodFrameContainer');

    const threeDSMethodData = {
        threeDSServerTransID: threeDSServerTransID,
        threeDSMethodNotificationURL: methodNotificationURL
    }
    const threeDSMethodDataBase64 = encode(JSON.stringify(threeDSMethodData))
    const methodIframeName = 'methodIframe'

    const html = document.createElement('html');
    const body = document.createElement('body');
    const methodIframe = createIframe(frameContainer, methodIframeName, methodIframeName, 0, 0)
    const methodForm = createForm('threeDSMethodForm', threeDSMethodURL, methodIframe.name)
    const input = createInput('threeDSMethodData', threeDSMethodDataBase64)

    methodForm.appendChild(input);
    body.appendChild(methodForm);
    html.appendChild(body);
    methodIframe.appendChild(html);

    methodForm.submit();
}

function SendChallengeRequest(acsURL, creq, sessionData) {
    let frameContainer = document.getElementById('challengeFrameContainer');

    const windowSize = getWindowSize(creq.challengeWindowSize)

    const creqBase64 = encode(JSON.stringify(creq));
    const sessionDataBase64 = encode(JSON.stringify(sessionData))
    const challengeIframeName = 'challengeIframe'

    const html = document.createElement('html');
    const body = document.createElement('body');
    const challengeIframe = createIframe(frameContainer, challengeIframeName, challengeIframeName, windowSize[0], windowSize[1])
    const form = createForm('threeDSCReqForm', acsURL, challengeIframe.name)
    const creqInput = createInput('creq', creqBase64)
    const sessionDataInput = createInput('threeDSSessionData', sessionDataBase64)

    form.appendChild(creqInput);
    form.appendChild(sessionDataInput);
    body.appendChild(form);
    html.appendChild(body);
    challengeIframe.appendChild(html);

    form.submit();
}

const getWindowSize = (challengeWindowSize = '05') => {
    switch (challengeWindowSize) {
        case '01':
            return ['250px', '400px'];
        case '02':
            return ['390px', '400px'];
        case '03':
            return ['500px', '600px'];
        case '04':
            return ['600px', '400px'];
        case '05':
            return ['100%', '100%'];
        default:
            throw Error(`Selected window size ${challengeWindowSize} is not supported`);
    }
};

function createIframe(container, name, id, width = '0', height = '0', onLoadCallback) {
    if (!container || !name || !id) {
        throw Error('Not all required fields have value');
    }
    if (!(container instanceof HTMLElement)) {
        throw Error('Container must be a HTML element');
    }

    const iframe = document.createElement('iframe');
    iframe.name = name;
    iframe.width = width;
    iframe.height = height;
    iframe.setAttribute('id', id);
    iframe.setAttribute('frameborder', '0');
    iframe.setAttribute('border', '0');

    if (onLoadCallback && typeof onLoadCallback === 'function') {
        if (iframe.attachEvent) {
            iframe.attachEvent('onload', onLoadCallback);
        } else {
            iframe.onload = onLoadCallback;
        }
    }

    container.appendChild(iframe);
    return iframe;
}

function createForm(formName, formAction, formTarget) {
    const form = document.createElement('form');
    form.name = formName;
    form.action = formAction;
    form.method = 'POST';
    form.target = formTarget;
    return form
}

function createInput(inputName, inputValue) {
    const input = document.createElement('input');
    input.name = inputName;
    input.value = inputValue;
    return input
}

function encode (str) {
    return btoa(str).replace(/\+/g, '-')
        .replace(/\//g, '_')
        .replace(/=/g, '')
}