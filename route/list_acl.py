from .tool.func import *

def list_acl_2(conn):
    curs = conn.cursor()

    div = '<ul>'

    curs.execute(db_change("select title, decu, dis, view, why from acl where decu != '' or dis != '' or view != '' order by title desc"))
    list_data = curs.fetchall()
    for data in list_data:
        if not re.search('^user:', data[0]) and not re.search('^file:', data[0]):
            acl = []
            for i in range(1, 4):
                if data[i] == 'admin':
                    acl += [load_lang('admin')]
                elif data[i] == 'user':
                    acl += [load_lang('member')]
                elif data[i] == '':
                    acl += [load_lang('normal')]
                else:
                    acl += [data[i]]

            curs.execute(db_change("select time from re_admin where what like ? order by time desc limit 1"), ['acl (' + data[0] + ')%'])
            time_data = curs.fetchall()
            if time_data:
                time_data = ' | ' + time_data[0][0]
            else:
                time_data = ''

            div += '' + \
                '<li>' + \
                    '<a href="/w/' + url_pas(data[0]) + '">' + data[0] + '</a> | ' + \
                    load_lang('document_acl') + ' : ' + acl[0] + ' | ' + \
                    load_lang('discussion_acl') + ' : ' + acl[1] + ' | ' + \
                    load_lang('view_acl') + ' : ' + acl[2] + \
                    time_data + \
                '</li>' + \
            ''

    div += '</ul>'

    return easy_minify(flask.render_template(skin_check(),
        imp = [load_lang('acl_document_list'), wiki_set(), custom(), other2([0, 0])],
        data = div,
        menu = [['other', load_lang('return')]]
    ))