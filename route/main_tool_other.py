from .tool.func import *

async def main_tool_other():
    with get_db_connect() as conn:
        return easy_minify(conn, flask.render_template(skin_check(conn),
            imp = [get_lang(conn, 'other_tool'), await wiki_set(), await wiki_custom(conn), wiki_css([0, 0])],
            data = render_simple_set(conn, '''
                <h2>''' + get_lang(conn, 'user_tool') + '''</h2>
                <ul>
                    <li><a href="/manager/6">''' + get_lang(conn, 'user_tool') + '''</a></li>
                </ul>
                <h2>''' + get_lang(conn, 'list') + '''</h2>
                <h3>''' + get_lang(conn, 'admin') + '''</h3>
                <ul>               
                    <li><a href="/list/admin">''' + get_lang(conn, 'admin_list') + '''</a></li>
                    <li><a href="/list/admin/auth_use">''' + get_lang(conn, 'authority_use_list') + '''</a></li>
                </ul>
                <h3>''' + get_lang(conn, 'discussion') + '''</h3>
                <ul>
                    <li><a href="/recent_discuss">''' + get_lang(conn, 'recent_discussion') + '''</a></li>
                </ul>
                <h3>''' + get_lang(conn, 'document') + '''</h3>
                <ul>
                    <li><a href="/recent_change">''' + get_lang(conn, 'recent_change') + '''</a></li>
                    <li><a href="/list/document/all">''' + get_lang(conn, 'all_document_list') + '''</a></li>
                    <li><a href="/list/document/acl">''' + get_lang(conn, 'acl_document_list') + '''</a></li>
                    <li><a href="/list/document/need">''' + get_lang(conn, 'need_document') + '''</a></li>
                    <li><a href="/list/document/long">''' + get_lang(conn, 'long_page') + '''</a></li>
                    <li><a href="/list/document/short">''' + get_lang(conn, 'short_page') + '''</a></li>
                    <li><a href="/list/document/old">''' + get_lang(conn, 'old_page') + '''</a></li>
                    <li><a href="/list/document/new">''' + get_lang(conn, 'new_page') + '''</a></li>
                    <li><a href="/list/document/no_link">''' + get_lang(conn, 'no_link_document_list') + '''</a></li>
                </ul>
                <h3>''' + get_lang(conn, 'user') + '''</h3>
                <ul>
                    <li><a href="/recent_block">''' + get_lang(conn, 'recent_ban') + '''</a></li>
                    <li><a href="/list/user">''' + get_lang(conn, 'member_list') + '''</a></li>
                </ul>
                <h3>''' + get_lang(conn, 'other') + '''</h3>
                <ul>
                    <li><a href="/list/file">''' + get_lang(conn, 'image_file_list') + '''</a></li>
                    <li><a href="/vote">''' + get_lang(conn, 'vote_list') + '''</a></li>
                    <li><a href="/bbs/main">''' + get_lang(conn, 'bbs_main') + '''</a></li>
                </ul>
                <h2>''' + get_lang(conn, 'other') + '''</h2>
                <ul>
                    <li><a href="/upload">''' + get_lang(conn, 'upload') + '''</a></li>
                    <li><a href="/manager/10">''' + get_lang(conn, 'search') + '''</a></li>
                </ul>
                <h2>''' + get_lang(conn, 'admin') + '''</h2>
                <ul>
                    <li><a href="/manager/1">''' + get_lang(conn, 'admin_tool') + '''</a></li>
                </ul>
            '''),
            menu = 0
        ))