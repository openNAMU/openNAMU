from .tool.func import *

def user_watch_list(tool):
    with get_db_connect() as conn:
        curs = conn.cursor()

        if tool == 'watch_list':
            div = load_lang("msg_whatchlist_lmt") + ' : 10 <hr class="main_hr">'
        else:
            div = ''

        ip = ip_check()

        if ip_or_user(ip) != 0:
            return redirect('/login')

        if tool == 'watch_list':
            curs.execute(db_change("select title from scan where type = '' and user = ?"), [ip])

            title_name = load_lang('watchlist')
        else:
            curs.execute(db_change("select title from scan where type = 'star' and user = ?"), [ip])

            title_name = load_lang('star_doc')

        data = curs.fetchall()
        for data_list in data:
            if tool == 'star_doc':
                curs.execute(db_change("select date from history where title = ? order by id + 0 desc limit 1"), [data_list[0]])
                get_data = curs.fetchall()
                if get_data:
                    plus = '(' + get_data[0][0] + ') '
                else:
                    plus = ''
            else:
                plus = ''

            div += '' + \
                '<li>' + \
                    '<a href="/w/' + url_pas(data_list[0]) + '">' + html.escape(data_list[0]) + '</a> ' + \
                    plus + \
                    '<a href="/' + ('star_doc' if tool == 'star_doc' else 'watch_list') + '/' + url_pas(data_list[0]) + '">(' + load_lang('delete') + ')</a>' + \
                '</li>' + \
            ''

        if data:
            div = '' + \
                '<ul class="opennamu_ul">' + div + '</ul>' + \
                '<hr class="main_hr">' + \
            ''

        div += '<a href="/manager/' + ('13' if tool == 'watch_list' else '16') + '">(' + load_lang('add') + ')</a>'

        return easy_minify(flask.render_template(skin_check(),
            imp = [title_name, wiki_set(), wiki_custom(), wiki_css([0, 0])],
            data = div,
            menu = [['user', load_lang('return')]]
        ))
