from .tool.func import *

async def main_tool_admin():
    return easy_minify(flask.render_template(await skin_check(),
        imp = [await get_lang('admin_tool'), await wiki_set(), await wiki_custom(), wiki_css([0, 0])],
        data = await render_simple_set('''
            <h2>''' + await get_lang('admin') + '''</h2>
            <ul>
                <li><a href="/manager/2">''' + await get_lang('document_setting') + '''</a></li>
                <li><a href="/acl_multiple">''' + await get_lang('mutiple_document_setting') + '''</a></li>
                <li><a href="/manager/3">''' + await get_lang('check_user') + '''</a></li>
                <li><a href="/auth/ban">''' + await get_lang('ban') + '''</a></li>
                <li><a href="/auth/ban/multiple">''' + await get_lang('multiple_ban') + '''</a></li>
                <li><a href="/manager/5">''' + await get_lang('authorize') + '''</a></li>
                <li><a href="/auth/give">''' + await get_lang('multiple_authorize') + '''</a></li>
                <li><a href="/auth/give_total">''' + await get_lang('auth_to_auth') + '''</a></li>
                <li><a href="/delete_multiple">''' + await get_lang('many_delete') + '''</a></li>
                <li><a href="/app_submit">''' + await get_lang('application_list') + '''</a></li>
            </ul>
            <h2>''' + await get_lang('owner') + '''</h2>
            <ul>
                <li><a href="/auth/list">''' + await get_lang('admin_group_list') + '''</a></li>
                <li><a href="/register">''' + await get_lang('add_user') + '''</a></li>
                <li><a href="/setting">''' + await get_lang('setting') + '''</a></li>
                <li><a href="/manager/18">''' + await get_lang('user_fix') + '''</a></li>
            </ul>
            <h3>''' + await get_lang('filter') + '''</h3>
            <ul>
                <li><a href="/filter/edit_filter">''' + await get_lang('edit_filter_list') + '''</a></li>
                <li><a href="/filter/inter_wiki">''' + await get_lang('interwiki_list') + '''</a></li>
                <li><a href="/filter/edit_top">''' + await get_lang('edit_tool_list') + '''</a></li>
                <li><a href="/filter/image_license">''' + await get_lang('image_license_list') + '''</a></li>
                <li><a href="/filter/email_filter">''' + await get_lang('email_filter_list') + '''</a></li>
                <li><a href="/filter/name_filter">''' + await get_lang('id_filter_list') + '''</a></li>
                <li><a href="/filter/file_filter">''' + await get_lang('file_filter_list') + '''</a></li>
                <li><a href="/filter/extension_filter">''' + await get_lang('extension_filter_list') + '''</a></li>
                <li><a href="/filter/document">''' + await get_lang('document_filter_list') + '''</a></li>
                <li><a href="/filter/outer_link">''' + await get_lang('outer_link_filter_list') + '''</a> (''' + await get_lang('beta') + ''')
                <li><a href="/filter/template">''' + await get_lang('template_document_list') + '''</a> (''' + await get_lang('beta') + ''')
            </ul>
            <h3>''' + await get_lang('server') + '''</h2>
            <ul>
                <li><a href="/restart">''' + await get_lang('wiki_restart') + '''</a></li>
                <li><a href="/shutdown">''' + await get_lang('wiki_shutdown') + '''</a></li>
                <li><a href="/update">''' + await get_lang('update') + '''</a></li>
            </ul>
            <h2>''' + await get_lang('version') + '''</h2>
            <ul>
                <li id="ver_send_2">''' + await get_lang('version') + ''' : </li>
                <li id="ver_send">''' + await get_lang('lastest') + ''' : </li>
            </ul>
            <h3>''' + await get_lang('skin_info') + '''</h3>
            <ul>
                <li><a href="/api/skin_info?all=true">''' + await get_lang('skin_info') + '''</a></li>
                <div id="ver_send_3"></div>
            </ul>
            <!-- JS : opennamu_do_insert_version -->
            <!-- JS : opennamu_do_insert_version_skin -->
        '''),
        menu = [['other', await get_lang('return')]]
    ))