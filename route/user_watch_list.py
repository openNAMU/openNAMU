from .tool.func import *

async def user_watch_list(tool):
    with get_db_connect() as conn:
        curs = conn.cursor()

        if tool == 'watch_list':
            div = get_lang(conn, "msg_whatchlist_lmt") + ' : 10 <hr class="main_hr">'
        else:
            div = ''

        ip = ip_check()

        if ip_or_user(ip) != 0:
            return redirect(conn, '/login')

        if tool == 'watch_list':
            curs.execute(db_change("select data from user_set where name = 'watchlist' and id = ?"), [ip])
            title_name = get_lang(conn, 'watchlist')
        else:
            curs.execute(db_change("select data from user_set where name = 'star_doc' and id = ?"), [ip])
            title_name = get_lang(conn, 'star_doc')

        data = curs.fetchall()
        for data_list in data:
            curs.execute(db_change("select date from history where title = ? order by id + 0 desc limit 1"), [data_list[0]])
            get_data = curs.fetchall()
            plus = '(' + get_data[0][0] + ') ' if get_data else ''
            
            div += '' + \
                '<li>' + \
                    '<a href="/w/' + url_pas(data_list[0]) + '">' + html.escape(data_list[0]) + '</a> ' + \
                    plus + \
                    '<a href="/' + ('star_doc' if tool == 'star_doc' else 'watch_list') + '/' + url_pas(data_list[0]) + '">(' + get_lang(conn, 'delete') + ')</a>' + \
                '</li>' + \
            ''

        if data:
            div = '' + \
                '<ul>' + div + '</ul>' + \
                '<hr class="main_hr">' + \
            ''

        div += '<a href="/manager/' + ('13' if tool == 'watch_list' else '16') + '">(' + get_lang(conn, 'add') + ')</a>'

        return easy_minify(conn, flask.render_template(skin_check(conn),
            imp = [title_name, await wiki_set(), await wiki_custom(conn), wiki_css([0, 0])],
            data = div,
            menu = [['user', get_lang(conn, 'return')]]
        ))
