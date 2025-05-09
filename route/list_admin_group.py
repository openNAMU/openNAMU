from .tool.func import *

async def list_admin_group():
    with get_db_connect() as conn:
        curs = conn.cursor()

        list_data = '<ul>'
        org_acl_list = get_default_admin_group()

        curs.execute(db_change("select distinct name from alist order by name asc"))
        for data in curs.fetchall():
            if await acl_check('', 'owner_auth', '', '') != 1 and not data[0] in org_acl_list:
                delete_admin_group = ' <a href="/auth/list/delete/' + url_pas(data[0]) + '">(' + get_lang(conn, "delete") + ')</a>'
            else:
                delete_admin_group = ''

            list_data += '' + \
                '<li>' + \
                    '<a href="/auth/list/add/' + url_pas(data[0]) + '">' + html.escape(data[0]) + '</a>' + \
                    delete_admin_group + \
                '</li>' + \
            ''

        list_data += '' + \
            '</ul>' + \
            '<hr class="main_hr">' + \
            '<a href="/manager/8">(' + get_lang(conn, 'add') + ')</a>' + \
        ''

        return easy_minify(conn, flask.render_template(skin_check(conn),
            imp = [get_lang(conn, 'admin_group_list'), await wiki_set(), await wiki_custom(conn), wiki_css([0, 0])],
            data = list_data,
            menu = [['manager', get_lang(conn, 'return')]]
        ))