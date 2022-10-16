const url = 'https://validator.stpr.cc';
const tid = document.getElementById("grjs").getAttribute('tid'); // Get tracker id

fetch(url+'/validator?tracker_uuid='+tid+"&url="+escape(window.location.href)) // do request to validator
.then(r => {
    if (!r.ok) {
    throw new Error(`failed ${r.status}`)
    }
    return r.text()
})
.then(d => {
    // Here we can do sometring if request is from gov ip.
})
.catch(e => console.log(e))
