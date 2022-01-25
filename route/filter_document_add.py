from .tool.func import *

def filter_document_add(name = ''):
    with get_db_connect() as conn:
        curs = conn.cursor()

        if flask.request.method == 'POST':
            if admin_check(None, 'plus_document_filter') != 1:
                return re_error('/error/3')

            post_name = flask.request.form.get('name', '')
            if post_name == '':
                return redirect('/filter/document/list')
            
            post_acl = flask.request.form.get('acl', '')
            post_regex = flask.request.form.get('regex', '')
            try:
                re.compile(post_regex)
            except:
                return re_error('/error/23')

            curs.execute(db_change('insert into html_filter (html, kind, plus, plus_t) values (?, ?, ?, ?)'), [
                post_name,
                'document',
                post_regex,
                post_acl
            ])

            conn.commit()

            return redirect('/filter/document/list')
        else:
            stat = 'disabled' if admin_check() != 1 else ''
            acl_list = get_acl_list()
            
            curs.execute(db_change("select plus, plus_t from html_filter where html = ? and kind = 'document'"), [name])
            db_data = curs.fetchall()
            acl_list = [['selected' if db_data and db_data[0][1] == i else '', i] for i in acl_list]

            return easy_minify(flask.render_template(skin_check(),
                imp = [load_lang('document_filter_add'), wiki_set(), wiki_custom(), wiki_css([0, 0])],
                data =  '''
                        <form method="post">
                            ''' + load_lang('name') + '''
                            <hr class="main_hr">
                            <input value="''' + html.escape(name) + '''" type="text" name="name">
                            <hr class="main_hr">
                            ''' + load_lang('regex') + '''
                            <hr class="main_hr">
                            <input value="''' + (html.escape(db_data[0][0]) if db_data else '') + '''" type="text" name="regex">
                            <hr class="main_hr">
                            <a href="/acl/Test#exp">''' + load_lang('acl') + '''</a>
                            <hr class="main_hr">
                            <select name="acl">
                                ''' + ''.join(['<option ' + i[0] + ' value=' + i[1] + '>' + ('normal' if i[1] == '' else i[1]) + '</option>' for i in acl_list]) + '''
                            </select>
                            <hr class="main_hr">
                            <button ''' + stat + ''' type="submit">''' + load_lang('add') + '''</button>
                        </form>
                        ''',
                menu = [['filter/document/list', load_lang('return')]]
            ))