async function postData(url, data) {
    const response = await fetch(url, {
        method: 'POST',
        cache: 'no-cache',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify(data)
    });
    if (response.ok) {
        return await response.json();
    } else {
        throw await response.json();
    }
}

async function getData(url) {
    const response = await fetch(url, {
        method: 'GET',
        cache: 'no-cache',
        headers: {
            'Content-Type': 'application/json'
        },
    })
    if (response.ok) {
        return await response.json();
    } else {
        throw await response.json();
    }
}

const EventHandling = {
    data() {
        return {
            // variables related to shorten api
            original_url: '',
            short_url_view: '',
            shorten_error: false,
            // variables related to original url lookup api
            short_url: '',
            original_url_view: '',
            lookup_error: false,
        }
    },
    methods: {
        genShortUrl() {
            postData('http://localhost:8080/api/shorten', {
                url: this.original_url
            }).then(data => {
                this.short_url_view = data.link
                this.shorten_error = false
            }).catch(error => {
                console.error(error)
                this.short_url_view = error.message
                this.shorten_error = true
            })
        },
        getOriginalUrl() {
            if (this.short_url.length === 0) {
                console.log('skip')
                return
            }
            getData(`http://localhost:8080/api/url/${this.short_url}`)
                .then(data => {
                    this.original_url_view = data.url
                    this.lookup_error = true
                }).catch(error => {
                    console.error(error)
                    this.original_url_view = error.message
                    this.lookup_error = true
                })
        }
    }
}

Vue.createApp(EventHandling).mount('#url-shortener')