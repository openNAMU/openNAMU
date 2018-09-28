<div id="plus">
</div>
<script>
    function addZero(i) {
        if(i < 10) {
            i = "0" + i;
        }
        return i;
    }

    setInterval(
        function() {
            var d = new Date();
            var date = d.getFullYear() + '-' + addZero(d.getMonth() + 1) + '-' + d.getDate() + ' ' + addZero(d.getHours()) + ':' + addZero(d.getMinutes()) + ':' + addZero(d.getSeconds());

            var url = "/api/topic/''' + name + '''/sub/''' + sub + '''";

            var xhr = new XMLHttpRequest();
            
            xhr.open("GET", url);
            xhr.send();

            xhr.onreadystatechange = function() {
                if(xhr.status == 200) {
                    var data = JSON.parse(xhr.responseText);

                    data.forEach(function(element) {
                        document.getElementById('plus').value += element['id'] + element['data'];
                    });
                }
            }
        }
    , 3000);
</script>