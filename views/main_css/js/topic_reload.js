function topic_reload() {
    location.href.endsWith('#reload') ? location.reload(true) : location.href = '#reload';
}