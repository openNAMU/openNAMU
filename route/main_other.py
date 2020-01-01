from .tool.func import *

def main_other_2(conn):
    curs = conn.cursor()

    return easy_minify(flask.render_template(skin_check(),
        imp = [load_lang('other_tool'), wiki_set(), custom(), other2([0, 0])],
        data = '''
            <h2>''' + load_lang('record') + '''</h2>
            <ul>
                <li><a href="/manager/6">''' + load_lang('edit_record') + '''</a></li>
                <li><a href="/manager/7">''' + load_lang('discussion_record') + '''</a></li>
            </ul>
            <br>
            <h2>''' + load_lang('list') + '''</h2>
            <ul>
                <li><a href="/admin_list">''' + load_lang('admin_list') + '''</a></li>
                <li><a href="/not_close_topic">''' + load_lang('open_discussion_list') + '''</a></li>
                <li><a href="/title_index">''' + load_lang('all_document_list') + '''</a></li>
                <li><a href="/acl_list">''' + load_lang('acl_document_list') + '''</a></li>
                <li><a href="/please">''' + load_lang('need_document') + '''</a></li>
                <li><a href="/block_log">''' + load_lang('recent_ban') + '''</a></li>
                <li><a href="/user_log">''' + load_lang('member_list') + '''</a></li>
                <li><a href="/admin_log">''' + load_lang('authority_use_list') + '''</a></li>
                <li><a href="/old_page">''' + load_lang('old_page') + '''</a></li>
            </ul>
            <br>
            <h2>''' + load_lang('other') + '''</h2>
            <ul>
                <li><a href="/upload">''' + load_lang('upload') + '''</a></li>
                <li><a href="/manager/10">''' + load_lang('search') + '''</a></li>
            </ul>
            <br>
            <h2>''' + load_lang('admin') + '''</h2>
            <ul>
                <li><a href="/manager/1">''' + load_lang('admin_tool') + '''</a></li>
            </ul>
        ''',
        menu = 0
    ))