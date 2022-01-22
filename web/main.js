async function postData(url, data) {
    const response = await fetch(url, {
        method: 'POST',
        cache: 'no-cache',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify(data)
    });
    return response.json();
}

async function getData(url) {
    const response = await fetch(url, {
        method: 'GET',
        cache: 'no-cache',
        headers: {
            'Content-Type': 'application/json'
        },
    })
    return response.json();
}

const EventHandling = {
    data() {
        return {
            original_url: 'Hello Vue.js!',
            short_url_view: '',
            short_url: 'short url',
            original_url_view: '',
        }
    },
    methods: {
        genShortUrl() {
            postData('http://localhost:8080/api/shorten', {
                url: this.original_url
            })
            .then(data => {
                this.short_url_view = data.link
            })
        },
        getOriginalUrl() {
            getData(`http://localhost:8080/api/url/${this.short_url}`)
            .then(data => {
                this.original_url_view = data.url
            })
        }
    }
}

Vue.createApp(EventHandling).mount('#url-shortener')