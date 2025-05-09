from .tool.func import *

from .go_api_list_old_page import api_list_old_page

async def list_old_page(num = 1, set_type = 'old'):
    with get_db_connect() as conn:
        title = ''
        if set_type == 'old':
            title = get_lang(conn, 'old_page')
        else:
            title = get_lang(conn, 'new_page')

        data = await api_list_old_page(num, set_type)
        data = data["data"]

        data_html = ''

        for for_a in range(len(data)):
            doc_name_encoded = url_pas(data[for_a][0])
            doc_title_filtered = html.escape(data[for_a][0])

            right = f'<a href="/w/{doc_name_encoded}">{doc_title_filtered}</a> '

            data_html += await opennamu_make_list(right, data[for_a][1])

        data_html += get_next_page_bottom(conn, f'/list/document/{set_type}/{{}}', int(num), data)

        return easy_minify(conn, flask.render_template(skin_check(conn),
            imp = [title, await wiki_set(), await wiki_custom(conn), wiki_css([0, 0])],
            data = data_html,
            menu = [['other', get_lang(conn, 'return')]]
        ))