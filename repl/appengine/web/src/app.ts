/**
 * @fileoverview basic javascript handling repl ui
 */


 (function() {

  interface Stringer {
    toString() : String
  }
  
  const endpoint = '/api';
  const newInput = `<div class="statement-block">
      &gt; <input class="repl-stmt repl-stmt-new" value=""><br>
      <code class="repl-out-ok"></code>
    </div>`;
  
  
  function trimLen(str : String) {
    if (str.length > 80) {
      str = str.slice(0, 77) + '...';
    }
    return str;
  }
  
  function quoteattr(str : String) {
    // TODO find a supported library to do this...
    return str
        .replace(/&/g, '&amp;')
        .replace(/'/g, '&apos;')
        .replace(/"/g, '&quot;')
        .replace(/</g, '&lt;')
        .replace(/>/g, '&gt;');
  }

  function showError(error: Stringer) {
    document.querySelector('pre#error-box')
        .innerHTML = quoteattr(error.toString())
  }

  function render (request, response) {
    let insertionPoint = document.querySelector('div.input-block');
    let html = '';
    let lastResp = '';
  
    request.commands.forEach(function(reqLine: String, i: number) {
      let resp : any = null;
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
  
    html += newInput;
  
    insertionPoint.innerHTML = html;
  
    document.querySelector('code#result').innerHTML = quoteattr(lastResp);
    readyForInput();

  }

  function submit() {
    let inputs = document.querySelectorAll('input.repl-stmt');
    let req = {'commands': []};
    inputs.forEach(function(el) {
      let inp = <HTMLInputElement>el;
      if (inp.value != null && inp.value.trim() != '') {
        req.commands.push(inp.value);
      }
    });
    let payload = JSON.stringify(req);
    console.log('submitting: %s', payload);
    window
        .fetch(endpoint, {
          method: 'POST',
          headers: {'Content-Type': 'application/json'},
          body: payload
        })
        .then(response => response.json())
        .then(data => render(req, data))
        .catch(ex => showError(ex));
  }

  function readyForInput() {
    let mainInput = <HTMLInputElement>document.querySelector('input.repl-stmt-new');
    mainInput.focus();

    mainInput.addEventListener('keyup', function(e : KeyboardEvent) {
      if (e.code == 'Enter' && !e.metaKey && !e.ctrlKey) {
        e.preventDefault();
        submit();
      }
    });
  }
  
  window.addEventListener('load', () => {
    let evalBtn = document.querySelector('button#evaluate');
    let addBtn = document.querySelector('button#add-statement');
    if (!evalBtn || !addBtn) {
      console.log('Unable to bind control callbacks.');
      return;
    }
    let insertionPoint = document.querySelector('div.input-block');
    insertionPoint.innerHTML = newInput;
    evalBtn.addEventListener('click', () => submit());
    readyForInput();
  });

 })();