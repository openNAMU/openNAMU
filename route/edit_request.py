from .tool.func import *

from .view_diff import view_diff_do

async def edit_request(name = 'Test', do_type = ''):
    with get_db_connect() as conn:
        curs = conn.cursor()

        disabled = ""
        if await acl_check(name, 'document_edit') == 1:
            disabled = "disabled"

        curs.execute(db_change("select id from history where title = ? order by id + 0 desc"), [name])
        doc_ver = curs.fetchall()
        doc_ver = doc_ver[0][0] if doc_ver else '0'

        if doc_ver == '0':
            if await acl_check(name, 'document_make_acl') == 1:
                disabled = "disabled"

        curs.execute(db_change("select set_data from data_set where doc_name = ? and doc_rev = ? and set_name = 'edit_request_data'"), [name, doc_ver])
        db_data = curs.fetchall()
        if not db_data:
            return redirect(conn, '/edit/' + url_pas(name))
        
        edit_request_data = db_data[0][0]

        curs.execute(db_change("select set_data from data_set where doc_name = ? and doc_rev = ? and set_name = 'edit_request_user'"), [name, doc_ver])
        db_data = curs.fetchall()
        edit_request_user = db_data[0][0] if db_data else ''

        curs.execute(db_change("select set_data from data_set where doc_name = ? and doc_rev = ? and set_name = 'edit_request_date'"), [name, doc_ver])
        db_data = curs.fetchall()
        edit_request_date = db_data[0][0] if db_data else ''

        curs.execute(db_change("select set_data from data_set where doc_name = ? and doc_rev = ? and set_name = 'edit_request_send'"), [name, doc_ver])
        db_data = curs.fetchall()
        edit_request_send = db_data[0][0] if db_data else ''

        curs.execute(db_change("select set_data from data_set where doc_name = ? and doc_rev = ? and set_name = 'edit_request_leng'"), [name, doc_ver])
        db_data = curs.fetchall()
        edit_request_leng = db_data[0][0] if db_data else ''

        if flask.request.method == 'POST':
            if disabled != "":
                return redirect(conn, '/w/' + url_pas(name))
            
            curs.execute(db_change("select id from user_set where name = 'watchlist' and data = ?"), [name])
            for scan_user in curs.fetchall():
                await add_alarm(scan_user[0], edit_request_user, '<a href="/w/' + url_pas(name) + '">' + html.escape(name) + '</a>')

            if flask.request.form.get('check', '') == 'Y':
                curs.execute(db_change("delete from data where title = ?"), [name])
                curs.execute(db_change("insert into data (title, data) values (?, ?)"), [name, edit_request_data])
                        
                history_plus(conn, 
                    name,
                    edit_request_data,
                    edit_request_date,
                    edit_request_user,
                    edit_request_send,
                    edit_request_leng,
                    mode = 'edit_request'
                )
                
                render_set(conn, 
                    doc_name = name,
                    doc_data = edit_request_data,
                    data_type = 'backlink'
                )
            else:
                history_plus(conn, 
                    name,
                    edit_request_data,
                    edit_request_date,
                    edit_request_user,
                    edit_request_send,
                    '0',
                    mode = 'edit_request'
                )
                
            if do_type == 'from':
                return redirect(conn, '/edit/' + url_pas(name))
            else:
                return redirect(conn, '/w/' + url_pas(name))
        else:
            curs.execute(db_change("select data from data where title = ?"), [name])
            db_data = curs.fetchall()
            old_data = db_data[0][0] if db_data else ''

            result = view_diff_do(old_data, edit_request_data, 'r' + doc_ver, get_lang(conn, 'edit_request'))

            return easy_minify(conn, flask.render_template(skin_check(conn), 
                imp = [name, await wiki_set(), await wiki_custom(conn), wiki_css(['(' + get_lang(conn, 'edit_request_check') + ')', 0])],
                data = '''
                    <div id="opennamu_get_user_info">''' + html.escape(edit_request_user) + '''</div>
                    <hr class="main_hr">
                    ''' + edit_request_date + '''
                    <hr class="main_hr">
                    <input readonly value="''' + html.escape(edit_request_send) + '''">
                    <hr class="main_hr">
                    ''' + result + '''
                    <hr class="main_hr">
                    <form method="post">
                        <button ''' + disabled + ''' id="opennamu_save_button" type="submit" name="check" value="Y">''' + get_lang(conn, 'approve') + '''</button>
                        <button ''' + disabled + ''' id="opennamu_preview_button" type="submit" name="check" value="">''' + get_lang(conn, 'decline') + '''</button>
                        <hr class="main_hr">
                        <textarea readonly class="opennamu_textarea_500">''' + html.escape(edit_request_data) + '''</textarea>
                    </form>
                ''',
                menu = 0
            ))