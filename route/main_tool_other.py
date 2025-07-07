from .tool.func import *

async def main_tool_other():
    with get_db_connect() as conn:
        return easy_minify(flask.render_template(await skin_check(conn),
            imp = [await get_lang('other_tool'), await wiki_set(), await wiki_custom(conn), wiki_css([0, 0])],
            data = await render_simple_set('''
                <h2>''' + await get_lang('user_tool') + '''</h2>
                <ul>
                    <li><a href="/manager/6">''' + await get_lang('user_tool') + '''</a></li>
                </ul>
                <h2>''' + await get_lang('list') + '''</h2>
                <h3>''' + await get_lang('admin') + '''</h3>
                <ul>               
                    <li><a href="/list/admin">''' + await get_lang('admin_list') + '''</a></li>
                    <li><a href="/list/admin/auth_use">''' + await get_lang('authority_use_list') + '''</a></li>
                </ul>
                <h3>''' + await get_lang('discussion') + '''</h3>
                <ul>
                    <li><a href="/recent_discuss">''' + await get_lang('recent_discussion') + '''</a></li>
                </ul>
                <h3>''' + await get_lang('document') + '''</h3>
                <ul>
                    <li><a href="/recent_change">''' + await get_lang('recent_change') + '''</a></li>
                    <li><a href="/list/document/all">''' + await get_lang('all_document_list') + '''</a></li>
                    <li><a href="/list/document/acl">''' + await get_lang('acl_document_list') + '''</a></li>
                    <li><a href="/list/document/need">''' + await get_lang('need_document') + '''</a></li>
                    <li><a href="/list/document/long">''' + await get_lang('long_page') + '''</a></li>
                    <li><a href="/list/document/short">''' + await get_lang('short_page') + '''</a></li>
                    <li><a href="/list/document/old">''' + await get_lang('old_page') + '''</a></li>
                    <li><a href="/list/document/new">''' + await get_lang('new_page') + '''</a></li>
                    <li><a href="/list/document/no_link">''' + await get_lang('no_link_document_list') + '''</a></li>
                </ul>
                <h3>''' + await get_lang('user') + '''</h3>
                <ul>
                    <li><a href="/recent_block">''' + await get_lang('recent_ban') + '''</a></li>
                    <li><a href="/list/user">''' + await get_lang('member_list') + '''</a></li>
                </ul>
                <h3>''' + await get_lang('other') + '''</h3>
                <ul>
                    <li><a href="/list/file">''' + await get_lang('image_file_list') + '''</a></li>
                    <li><a href="/vote">''' + await get_lang('vote_list') + '''</a></li>
                    <li><a href="/bbs/main">''' + await get_lang('bbs_main') + '''</a></li>
                </ul>
                <h2>''' + await get_lang('other') + '''</h2>
                <ul>
                    <li><a href="/upload">''' + await get_lang('upload') + '''</a></li>
                    <li><a href="/manager/10">''' + await get_lang('search') + '''</a></li>
                </ul>
                <h2>''' + await get_lang('admin') + '''</h2>
                <ul>
                    <li><a href="/manager/1">''' + await get_lang('admin_tool') + '''</a></li>
                </ul>
            '''),
            menu = 0
        ))