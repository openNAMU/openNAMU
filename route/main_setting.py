from .tool.func import *

async def main_setting():
    with get_db_connect() as conn:
        li_list = [
            ['main', await get_lang('main_setting')],
            ['phrase', await get_lang('text_setting')],
            ['robot', 'robots.txt'],
            ['external', await get_lang('ext_api_req_set')],
            ['head', await get_lang('main_head')],
            ['body/top', await get_lang('main_body')],
            ['body/bottom', await get_lang('main_bottom_body')],
            ['sitemap_set', await get_lang('sitemap_management')],
            ['top_menu', await get_lang('top_menu_setting')],
            ['skin_set', await get_lang('main_skin_set_default')],
            ['404_page', await get_lang('404_page_setting')]
        ]

        li_data = ''.join(['<li><a href="/setting/' + str(li[0]) + '">' + li[1] + '</a></li>' for li in li_list])

        return easy_minify(flask.render_template(await skin_check(),
            imp = [await get_lang('setting'), await wiki_set(), await wiki_custom(), wiki_css([0, 0])],
            data = '<h2>' + await get_lang('list') + '</h2><ul>' + li_data + '</ul>',
            menu = [['manager', await get_lang('return')]]
        ))