easter_egg = [
"""
<iframe width="560" height="315" src="https://www.youtube.com/embed/O6BJiije6m4" frameborder="0" allow="autoplay; encrypted-media" allowfullscreen></iframe>
<br>
<br>
<div>
    노라조와 영접하십셔.
</div>
""",
"""
<iframe width="560" height="315" src="https://www.youtube.com/embed/qnYWeY6aybg?autoplay=1&start=97" frameborder="0" allow="accelerometer; autoplay; encrypted-media; gyroscope; picture-in-picture" allowfullscreen></iframe>
<br>
<br>
<div>
    사람들은 알수 있어 누가 착한지 나쁜지
</div>
""",
"""
<iframe width="560" height="315" src="https://www.youtube.com/embed/oIa3BFBHJFI?autoplay=1&loop=1" frameborder="0" allow="accelerometer; autoplay; encrypted-media; gyroscope; picture-in-picture" allowfullscreen></iframe>
<br>
<br>
<div>
    @(^0^)==@태보는 지금 전 세계적으로 선풍적인 인기를 끌고 있는데요!!!!!!!!!!!@==(^0^)@
</div>
""",
"""<iframe width="560" height="315" src="https://www.youtube.com/embed/kxopViU98Xo?autoplay=1" frameborder="0" allow="accelerometer; autoplay; encrypted-media; gyroscope; picture-in-picture" allowfullscreen></iframe>""",
"""<iframe width="560" height="315" src="https://www.youtube.com/embed/dQw4w9WgXcQ?autoplay=1" frameborder="0" allow="accelerometer; autoplay; encrypted-media; gyroscope; picture-in-picture" allowfullscreen></iframe>"""
]

easter_egg_config_js = '''
<script>
document.getElementById('enable-easter-egg').addEventListener("change", function() {
    $.ajax({
        type: "POST",
        url: "/request/egg_config",
        data: {request: "config", data: document.getElementById('enable-easter-egg').checked},
        dataType: "json",
        success: function (data) {
            document.getElementById('easter-egg-status').innerHTML = '서버에 변경 사항이 저장되었습니다.';
        }
    });
});
</script>
'''