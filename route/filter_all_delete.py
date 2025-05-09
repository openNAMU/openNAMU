from .tool.func import *

async def filter_all_delete(tool, name = 'Test'):
    with get_db_connect() as conn:
        curs = conn.cursor()
        
        if await acl_check(tool = 'owner_auth', memo = 'del_' + tool) == 1:
            return await re_error(conn, 3)

        if tool == 'inter_wiki':
            curs.execute(db_change("delete from html_filter where html = ? and kind = 'inter_wiki'"), [name])
            curs.execute(db_change("delete from html_filter where html = ? and kind = 'inter_wiki_sub'"), [name])
        elif tool == 'edit_filter':
            curs.execute(db_change("delete from html_filter where html = ? and kind = 'regex_filter'"), [name])
        elif tool == 'name_filter':
            curs.execute(db_change("delete from html_filter where html = ? and kind = 'name'"), [name])
        elif tool == 'file_filter':
            curs.execute(db_change("delete from html_filter where html = ? and kind = 'file'"), [name])
        elif tool == 'email_filter':
            curs.execute(db_change("delete from html_filter where html = ? and kind = 'email'"), [name])
        elif tool == 'image_license':
            curs.execute(db_change("delete from html_filter where html = ? and kind = 'image_license'"), [name])
        elif tool == 'extension_filter':
            curs.execute(db_change("delete from html_filter where html = ? and kind = 'extension'"), [name])
        elif tool == 'document':
            curs.execute(db_change("delete from html_filter where html = ? and kind = 'document'"), [name])
        elif tool == 'outer_link':
            curs.execute(db_change("delete from html_filter where html = ? and kind = 'outer_link'"), [name])
        elif tool == 'template':
            curs.execute(db_change("delete from html_filter where html = ? and kind = 'template'"), [name])
        else:
            curs.execute(db_change("delete from html_filter where html = ? and kind = 'edit_top'"), [name])

        return redirect(conn, '/filter/' + tool)