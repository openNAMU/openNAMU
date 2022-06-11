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
        
        self.data_nowiki = {}
        self.data_backlink = []
        
        self.data_toc = ''
        self.data_footnote = ''
        self.data_category = ''

    def do_render_text(self):
        # <b>
        self.render_data = re.sub(
            r"&#x27;&#x27;&#x27;((?:(?!&#x27;&#x27;&#x27;).)+)&#x27;&#x27;&#x27;",
            '<b>\g<1></b>',
            self.render_data
        )
        # <i>
        self.render_data = re.sub(
            r"&#x27;&#x27;((?:(?!&#x27;&#x27;).)+)&#x27;&#x27;",
            '<i>\g<1></i>',
            self.render_data
        )
        # <u>
        self.render_data = re.sub(
            r"__((?:(?!__).)+)__",
            '<u>\g<1></u>',
            self.render_data
        )
        
        # <sup>
        self.render_data = re.sub(
            r"\^\^\^((?:(?!\^\^\^).)+)\^\^\^",
            '<sup>\g<1></sup>',
            self.render_data
        )
        # <sup> 2
        self.render_data = re.sub(
            r"\^\^((?:(?!\^\^).)+)\^\^",
            '<sup>\g<1></sup>',
            self.render_data
        )
        
        # <sub>
        self.render_data = re.sub(
            r",,,((?:(?!,,,).)+),,,",
            '<sub>\g<1></sub>',
            self.render_data
        )
        # <sub> 2
        self.render_data = re.sub(
            r",,((?:(?!,,).)+),,",
            '<sub>\g<1></sub>',
            self.render_data
        )
        
        # <s>
        self.render_data = re.sub(
            r"--((?:(?!--).)+)--",
            '<s>\g<1></s>',
            self.render_data
        )
        # <s> 2
        self.render_data = re.sub(
            r"~~((?:(?!~~).)+)~~",
            '<s>\g<1></s>',
            self.render_data
        )

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
        self.render_data = html.escape(self.doc_data)
        
        self.render_data_js = ''

        self.do_render_text()
        self.do_render_last()
        
        return [
            self.render_data, # HTML
            self.render_data_js, # JS
            [] # Other
        ]