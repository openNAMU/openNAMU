from .tool.func import *

def main_setting():
    with get_db_connect() as conn:
        li_list = [
            ['main', load_lang('main_setting')],
            ['phrase', load_lang('text_setting')],
            ['robot', 'robots.txt'],
            ['external', load_lang('ext_api_req_set')],
            ['head', load_lang('main_head')],
            ['body/top', load_lang('main_body')],
            ['body/bottom', load_lang('main_bottom_body')],
            ['sitemap_set', load_lang('sitemap_management')],
            ['top_menu', load_lang('top_menu_setting')],
            ['skin_set', load_lang('main_skin_set_default')],
            ['acl', load_lang('main_acl_setting')]
        ]

        li_data = ''.join(['<li><a href="/setting/' + str(li[0]) + '">' + li[1] + '</a></li>' for li in li_list])

        return easy_minify(flask.render_template(skin_check(),
            imp = [load_lang('setting'), wiki_set(), wiki_custom(), wiki_css([0, 0])],
            data = '<h2>' + load_lang('list') + '</h2><ul class="opennamu_ul">' + li_data + '</ul>',
            menu = [['manager', load_lang('return')]]
        ))