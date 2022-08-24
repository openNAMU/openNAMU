from .func_tool import *

class class_do_render_namumark:
    def __init__(
        self,
        curs,
        doc_name, 
        doc_data, 
        doc_include
    ):
        self.curs = curs
        
        self.doc_data = doc_data
        self.doc_name = doc_name
        self.doc_include = doc_include
        
        self.data_temp_storage = {}
        self.data_temp_storage_count = 0
        self.data_backlink = []
        
        self.data_toc = ''
        self.data_footnote = ''
        self.data_category = ''

        self.render_data = self.doc_data
        self.render_data = html.escape(self.render_data)
        self.render_data = '<back_br>\n' + self.render_data + '\n<front_br>'
        self.render_data_js = ''

    def get_tool_temp_storage(self, data_A = '', data_B = ''):
        self.data_temp_storage_count += 1
        data_name = 'opennamu_render_' + str(self.data_temp_storage_count)

        self.data_temp_storage[data_name] = data_A
        self.data_temp_storage['/' + data_name] = data_B

        return data_name

    def get_tool_data_restore(self, data):
        storage_count = self.data_temp_storage_count * 2
        storage_regex = r'<(\/?opennamu_render_(?:[0-9]+))>'

        while 1:
            if storage_count < 0:
                print('Error : render count overflow')

                break
            else:
                if re.search(storage_regex, data):
                    data = re.sub(
                        storage_regex, 
                        lambda match : self.data_temp_storage[match.group(1)], 
                        data,
                        1
                    )
                else:
                    break

            storage_count -= 1

        return data

    def do_render_text(self):
        # <b> function
        def do_render_text_bold(match):
            data = match.group(1)
            data_name = self.get_tool_temp_storage('<b>', '</b>')
            
            return '<' + data_name + '>' + data + '</' + data_name + '>'

        # <b>
        self.render_data = re.sub(
            r"&#x27;&#x27;&#x27;((?:(?!&#x27;&#x27;&#x27;).)+)&#x27;&#x27;&#x27;",
            do_render_text_bold,
            self.render_data
        )

        # <i> function
        def do_render_text_italic(match):
            data = match.group(1)
            data_name = self.get_tool_temp_storage('<i>', '</i>')
            
            return '<' + data_name + '>' + data + '</' + data_name + '>'

        # <i>
        self.render_data = re.sub(
            r"&#x27;&#x27;((?:(?!&#x27;&#x27;).)+)&#x27;&#x27;",
            do_render_text_italic,
            self.render_data
        )

        # <u> function
        def do_render_text_under(match):
            data = match.group(1)
            data_name = self.get_tool_temp_storage('<u>', '</u>')
            
            return '<' + data_name + '>' + data + '</' + data_name + '>'

        # <u>
        self.render_data = re.sub(
            r"__((?:(?!__).)+)__",
            do_render_text_under,
            self.render_data
        )
        
        # <sup> function
        def do_render_text_sup(match):
            data = match.group(1)
            data_name = self.get_tool_temp_storage('<sup>', '</sup>')
            
            return '<' + data_name + '>' + data + '</' + data_name + '>'

        # <sup>
        self.render_data = re.sub(
            r"\^\^\^((?:(?!\^\^\^).)+)\^\^\^",
            do_render_text_sup,
            self.render_data
        )
        # <sup> 2
        self.render_data = re.sub(
            r"\^\^((?:(?!\^\^).)+)\^\^",
            do_render_text_sup,
            self.render_data
        )

        # <sub> function
        def do_render_text_sub(match):
            data = match.group(1)
            data_name = self.get_tool_temp_storage('<sub>', '</sub>')
            
            return '<' + data_name + '>' + data + '</' + data_name + '>'
        
        # <sub>
        self.render_data = re.sub(
            r",,,((?:(?!,,,).)+),,,",
            do_render_text_sub,
            self.render_data
        )
        # <sub> 2
        self.render_data = re.sub(
            r",,((?:(?!,,).)+),,",
            do_render_text_sub,
            self.render_data
        )

        # <sub> function
        def do_render_text_strike(match):
            data = match.group(1)
            data_name = self.get_tool_temp_storage('<s>', '</s>')
            
            return '<' + data_name + '>' + data + '</' + data_name + '>'
        
        # <s>
        self.render_data = re.sub(
            r"--((?:(?!--).)+)--",
            do_render_text_strike,
            self.render_data
        )
        # <s> 2
        self.render_data = re.sub(
            r"~~((?:(?!~~).)+)~~",
            do_render_text_strike,
            self.render_data
        )
    
    def do_render_heading(self):
        pass

    def do_render_last(self):
        # remove front_br and back_br
        self.render_data = re.sub(
            r'\n<front_br>',
            '',
            self.render_data
        )
        self.render_data = re.sub(
            r'<back_br>\n',
            '',
            self.render_data
        )
        
        # \n to <br>
        self.render_data = re.sub(
            r'\n',
            '<br>',
            self.render_data
        )

    def __call__(self):
        self.do_render_text()
        self.do_render_heading()
        self.do_render_last()

        print(self.data_temp_storage)
        self.render_data = self.get_tool_data_restore(self.render_data)
        print('----')
        print(self.render_data)

        return [
            self.render_data, # HTML
            self.render_data_js, # JS
            [] # Other
        ]