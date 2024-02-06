from .tool.func import *

def main_tool_admin():
    with get_db_connect() as conn:
        return easy_minify(flask.render_template(skin_check(),
            imp = [load_lang('admin_tool'), wiki_set(), wiki_custom(), wiki_css([0, 0])],
            data = render_simple_set('''
                <h2>''' + load_lang('admin') + '''</h2>
                <ul class="opennamu_ul">
                    <li><a href="/manager/2">''' + load_lang('acl_change') + '''</a></li>
                    <li><a href="/manager/3">''' + load_lang('check_user') + '''</a></li>
                    <li><a href="/auth/give/ban">''' + load_lang('ban') + '''</a></li>
                    <li><a href="/auth/give/ban_multiple">''' + load_lang('multiple_ban') + '''</a></li>
                    <li><a href="/manager/17">''' + load_lang('release') + '''</a></li>
                    <li><a href="/manager/5">''' + load_lang('authorize') + '''</a></li>
                </ul>
                <h2>''' + load_lang('owner') + '''</h2>
                <ul class="opennamu_ul">
                    <li><a href="/auth/list">''' + load_lang('admin_group_list') + '''</a></li>
                    <li><a href="/delete_multiple">''' + load_lang('many_delete') + '''</a></li>
                    <li><a href="/app_submit">''' + load_lang('application_list') + '''</a></li>
                    <li><a href="/register">''' + load_lang('add_user') + '''</a></li>
                    <li><a href="/setting">''' + load_lang('setting') + '''</a></li>
                    <li><a href="/manager/18">''' + load_lang('user_fix') + '''</a></li>
                </ul>
                <h3>''' + load_lang('filter') + '''</h3>
                <ul class="opennamu_ul">
                    <li><a href="/filter/edit_filter">''' + load_lang('edit_filter_list') + '''</a></li>
                    <li><a href="/filter/inter_wiki">''' + load_lang('interwiki_list') + '''</a></li>
                    <li><a href="/filter/edit_top">''' + load_lang('edit_tool_list') + '''</a></li>
                    <li><a href="/filter/image_license">''' + load_lang('image_license_list') + '''</a></li>
                    <li><a href="/filter/email_filter">''' + load_lang('email_filter_list') + '''</a></li>
                    <li><a href="/filter/name_filter">''' + load_lang('id_filter_list') + '''</a></li>
                    <li><a href="/filter/file_filter">''' + load_lang('file_filter_list') + '''</a></li>
                    <li><a href="/filter/extension_filter">''' + load_lang('extension_filter_list') + '''</a></li>
                    <li><a href="/filter/document">''' + load_lang('document_filter_list') + '''</a></li>
                    <li><a href="/filter/outer_link">''' + load_lang('outer_link_filter_list') + '''</a> (''' + load_lang('beta') + ''')
                    <li><a href="/filter/template">''' + load_lang('template_document_list') + '''</a> (''' + load_lang('beta') + ''')
                </ul>
                <h3>''' + load_lang('server') + '''</h2>
                <ul class="opennamu_ul">
                    <li><a href="/restart">''' + load_lang('wiki_restart') + '''</a></li>
                    <li><a href="/shutdown">''' + load_lang('wiki_shutdown') + '''</a></li>
                    <li><a href="/update">''' + load_lang('update') + '''</a></li>
                </ul>
                <h2>''' + load_lang('version') + '''</h2>
                <ul class="opennamu_ul">
                    <li id="ver_send_2">''' + load_lang('version') + ''' : </li>
                    <li id="ver_send">''' + load_lang('lastest') + ''' : </li>
                </ul>
                <h3>''' + load_lang('skin_info') + '''</h3>
                <ul class="opennamu_ul">
                    <li><a href="/api/skin_info?all=true">''' + load_lang('skin_info') + '''</a></li>
                    <div id="ver_send_3"></div>
                </ul>
                <!-- JS : opennamu_do_insert_version -->
                <!-- JS : opennamu_do_insert_version_skin -->
            '''),
            menu = [['other', load_lang('return')]]
        ))