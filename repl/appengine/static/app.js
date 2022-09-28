/**
 * @fileoverview basic javascript handling repl ui
 * TODO(jdtatum): setup ts tool chain.
 */

 window.app = window.app || {};

 (function(app) {
 
 app.endpoint = '/api';
 
 app.newInput = `<div class="statement-block">
     &gt; <input class="repl-stmt repl-stmt-new" value=""><br>
     <code class="repl-out-ok"></code>
   </div>`;
 
 function trimLen(str) {
   if (str.length > 80) {
     str = str.substr(0, 77) + '...';
   }
   return str;
 }
 
 function quoteattr(str) {
   // TODO find a supported library to do this...
   return str
       .replace(/&/g, '&amp;')
       .replace(/'/g, '&apos;')
       .replace(/"/g, '&quot;')
       .replace(/</g, '&lt;')
       .replace(/>/g, '&gt;');
 }

 app.showError = function(error) {
    document.querySelector('pre#error-box')
        .innerHTML = quote(error.toString())
 };
  
 app.render = function(request, response) {
   let insertionPoint = document.querySelector('div.input-block');
   let html = '';
   let lastResp = '';
 
   request.commands.forEach(function(reqLine, i) {
     let resp = {};
     if (i < response.responses.length) {
       resp = response.responses[i];
     }
     console.log(response, resp);
     let error = !!resp.issue;
     let respLine = error ? resp.issue || 'err' : resp.replOutput || 'ok';
     lastResp = respLine;
     let preview = trimLen(respLine);
 
     html += `<div class="statement-block">
         &gt; <input class="repl-stmt" value="${quoteattr(reqLine)}"><br>
         <code class="${error ? 'repl-out-error' : 'repl-out-ok'}">${
         quoteattr(preview)}</code>
       </div>`;
   });
 
   html += app.newInput;
 
   insertionPoint.innerHTML = html;
 
   document.querySelector('code#result').innerHTML = quoteattr(lastResp);
   app.readyForInput();
 };
 
 app.submit = function() {
   let inputs = document.querySelectorAll('input.repl-stmt');
   let req = {'commands': []};
   inputs.forEach(function(inp) {
     if (inp.value != null && inp.value.trim() != '') {
       req.commands.push(inp.value);
     }
   });
   let payload = JSON.stringify(req);
   console.log('submitting: %s', payload);
   window
       .fetch(app.endpoint, {
         method: 'POST',
         headers: {'Content-Type': 'application/json'},
         body: payload
       })
       .then(response => response.json())
       .then(data => app.render(req, data))
       .catch(ex => console.log('TODO, ui for errors %s', ex));
 };
 
 app.readyForInput = function() {
   let mainInput = document.querySelector('input.repl-stmt-new');
   mainInput.focus();
   mainInput.addEventListener('keyup', function(e) {
     if (e.code == 'Enter' && !e.metaKey && !e.ctrlKey) {
       e.preventDefault();
       app.submit();
     }
   });
 };
 
 function setup() {
   let evalBtn = document.querySelector('button#evaluate');
   let addBtn = document.querySelector('button#add-statement');
   if (!evalBtn || !addBtn) {
     console.log('Unable to bind control callbacks.');
     return;
   }
   let insertionPoint = document.querySelector('div.input-block');
   insertionPoint.innerHTML = app.newInput;
   evalBtn.addEventListener('click', app.submit);
   addBtn.addEventListener('click', app.addInput);
   app.readyForInput();
 }
 
 window.addEventListener('load', setup);
 })(window.app);