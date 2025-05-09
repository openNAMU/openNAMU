from .tool.func import *

async def list_image_file(arg_num = 1, do_type = 0):
    with get_db_connect() as conn:
        curs = conn.cursor()

        sql_num = (arg_num * 50 - 50) if arg_num * 50 > 0 else 0

        list_data = ''
        if do_type == 0:
            list_data += '<a href="/list/image">(' + get_lang(conn, 'image') + ')</a>'
        else:
            list_data += '<a href="/list/file">(' + get_lang(conn, 'normal') + ')</a>'
        
        list_data += '<hr class="main_hr">'

        if do_type == 1:
            render_data = ''
            sub_data = ''
            count = 0

            curs.execute(db_change("select title from data where title like 'file:%' limit ?, 50"), [sql_num])
            data_list = curs.fetchall()
            for data in data_list:
                if count != 0 and count % 4 == 0:
                    render_data += '||\n'
                    render_data += sub_data + '||\n'
                    
                    sub_data = ''

                render_data += '|| [[' + data[0] + ']] '
                sub_data += '|| [[:' + data[0] + ']] '
                count += 1

            if render_data != '':
                render_data += '||\n'
                render_data += sub_data + '||'

            end_data = render_set(conn, 
                doc_name = '',
                doc_data = render_data,
                data_type = 'view',
                markup = 'namumark'
            )
            list_data += end_data
        else:
            list_data += '<ul>'

            curs.execute(db_change("select title from data where title like 'file:%' limit ?, 50"), [sql_num])
            data_list = curs.fetchall()
            for data in data_list:
                list_data += '<li><a href="/w/' + url_pas(data[0]) + '">' + html.escape(data[0]) + '</a></li>'

            list_data += '</ul>'

        if do_type == 0:
            list_data += get_next_page_bottom(conn, '/list/file/{}', arg_num, data_list)
        else:
            list_data += get_next_page_bottom(conn, '/list/image/{}', arg_num, data_list)

        return easy_minify(conn, flask.render_template(skin_check(conn),
            imp = [get_lang(conn, 'image_file_list'), await wiki_set(), await wiki_custom(conn), wiki_css([0, 0])],
            data = list_data,
            menu = [['other', get_lang(conn, 'return')]]
        ))