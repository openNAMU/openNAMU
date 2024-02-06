from .tool.func import *

def list_admin_group_2():
    with get_db_connect() as conn:
        curs = conn.cursor()

        list_data = '<ul class="opennamu_ul">'
        org_acl_list = get_default_admin_group()

        curs.execute(db_change("select distinct name from alist order by name asc"))
        for data in curs.fetchall():            
            if  admin_check() == 1 and \
                not data[0] in org_acl_list:
                delete_admin_group = ' <a href="/auth/list/delete/' + url_pas(data[0]) + '">(' + load_lang("delete") + ')</a>'
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
            '<a href="/manager/8">(' + load_lang('add') + ')</a>' + \
        ''

        return easy_minify(flask.render_template(skin_check(),
            imp = [load_lang('admin_group_list'), wiki_set(), wiki_custom(), wiki_css([0, 0])],
            data = list_data,
            menu = [['manager', load_lang('return')]]
        ))