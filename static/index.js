// This file is an example of the Javascript a merchant or PSP would need to add to
// their front-end in order to perform 3D Secure Authentication with Ravelin.

// Checkout calls the /checkout endpoint on the merchant backend initiating the checkout process.
function Checkout() {
    $('#payment').hide()
    $('#paymentProcessing').show()

    const requestBody = {
        accountNumber: document.getElementById('cardSelector').value,
    };

    console.log('Sending example merchant backend /checkout request using card ending in ' + requestBody.accountNumber.substr(-4))

    fetch(window.location.origin + '/checkout', {
        method: 'post',
        mode: 'cors',
        body: JSON.stringify(requestBody),
	    headers: {'Content-Type': 'application/json'}
    }).then(
        function (response) {
            if (response.status !== 200) {
                console.log('Looks like there was a problem. Status Code: ' + response.status);
                return;
            }

            console.log('/checkout response received')

            response.json().then(function (data) {
                if (data.threeDSMethodURL) {
                    console.log('threeDSMethodURL found, sending Method Request')
                    SendMethodRequest(data.threeDSMethodURL + '?success=true', data.threeDSServerTransID, window.location.origin + '/method-notification')
                    setTimeout(function() {
                        const msg = {
                            methodTimedOut: true,
                            threeDSServerTransID: data.threeDSServerTransID
                        };
                        window.postMessage(msg, "*");
                    }, 10000 )
                } else {
                    console.log('threeDSMethodURL not found, sending Authenticate Request')
                    Authenticate(data.threeDSServerTransID);
                }
            });
        }
    ).catch(function (err) {
        console.log('Error sending example merchant backend /checkout request', err);
    });
}

// Authenticate calls the /authenticate endpoint on the merchant backend, initiating the
// 3DS authentication process.
function Authenticate(threeDSServerTransID) {
    const requestBody = {
        productSKU: '10001',
        productQuantity: 1,
        accountNumber: document.getElementById('cardSelector').value,
        cardExpiryDate: '2205',
        threeDSServerTransID: threeDSServerTransID,
        browserData: GetBrowserData(),
    };

    console.log('Sending example merchant backend /authenticate request using card ending in ' + requestBody.accountNumber.substr(-4))

    fetch(window.location.origin + '/authenticate', {
        method: 'post',
        mode: 'cors',
        body: JSON.stringify(requestBody),
        headers: {'Content-Type': 'application/json'}
    }).then(
        function (response) {
            if (response.status !== 200) {
                console.log('Looks like there was a problem. Status Code: ' + response.status);
                return;
            }

            console.log('/authenticate response received')

            response.json().then(function (data) {
                if (data.error) {
                    console.log(data.error);
                    return
                }

                if (data.status === 'CHALLENGE_REQUIRED') {
                    console.log('Challenge Required, sending Challenge Request')

                    $('#payment').hide()
                    $('#paymentProcessing').hide()

                    const challengeRequest = {
                        messageType: 'CReq',
                        messageVersion: data.messageVersion,
                        threeDSServerTransID: data.threeDSServerTransID,
                        acsTransID: data.acsTransID,
                        challengeWindowSize: '03'
                    }
                    SendChallengeRequest(data.acsURL, challengeRequest, {})
                } else {
                    updatePage(data.status)
                }
            });
        }
    ).catch(function (err) {
        console.log('Error sending example merchant backend /authenticate request', err);
    });
}

// An event listener is used so that the Method and Challenge iframes can notify the
// parent page that processing has completed.
window.addEventListener('message', (e) => {
    if (e.origin === window.location.origin) {
        const event = e.data;

        // Method Notification
        if (event.hasOwnProperty('methodCompleted')) {
            console.log('Method Request completed');
            document.getElementById('methodIframe').remove();
            Authenticate(event.threeDSServerTransID);
        }
        
        if (event.hasOwnProperty('methodTimedOut')) {
            console.log('Method Request timed out');
            document.getElementById('methodIframe').remove();
            Authenticate(event.threeDSServerTransID);
        }

        // Challenge Notification
        if (event.hasOwnProperty('challengeCompleted')) {
            console.log('Challenge Request completed');
            document.getElementById('challengeIframe').remove()
            updatePage(event.status)
        }
    }
});

function updatePage(status) {
    $('#paymentProcessing').hide()
    if (status === 'SUCCESS') {
        $('#paymentSuccess').show()
    } else if (status === 'FAILED') {
        $('#paymentFailed').show()
    }
}

function resetPage() {
    $('#payment').show()
    $('#paymentProcessing').hide()
    $('#paymentSuccess').hide()
    $('#paymentFailed').hide()
}

function getTestCards() {
    fetch(window.location.origin + '/test-cards')
        .then(
            function (response) {
                if (response.status !== 200) {
                    console.warn('Failed to load test cards');
                    return;
                }
                response.json().then(function (data) {
                    let option;
                    if (data) {
                        let testCardSelect = document.getElementById('cardSelector');

                        for (let i = 0; i < data.length; i++) {
                            option = document.createElement('option');
                            option.text = data[i].description;
                            option.value = data[i].testPan;
                            testCardSelect.append(option);
                        }
                    }
                });
            }
        )
        .catch(function (err) {
            console.error('Failed to load test cards -', err);
        });
}

getTestCards()