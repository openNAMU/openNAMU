function render_markdown() {
    data = document.getElementById('render_contect').innerHTML;

    data = data.replace(/\*\*((?:(?!\*\*).)+)\*\*/, '<b>$1</b>');
    data = data.replace(/__((?:(?!__).)+)__/, '<b>$1</b>');
    data = data.replace(/\*((?:(?!\*).)+)\*/, '<i>$1</i>');
    data = data.replace(/_((?:(?!_).)+)_/, '<i>$1</i>');

    document.getElementById('render_contect').innerHTML = data;
}