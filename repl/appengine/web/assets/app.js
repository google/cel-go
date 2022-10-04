/**
 * @fileoverview basic javascript handling repl ui
 */
(function () {
    var endpoint = '/api';
    var newInput = "<div class=\"statement-block\">\n      &gt; <input class=\"repl-stmt repl-stmt-new\" value=\"\"><br>\n      <code class=\"repl-out-ok\"></code>\n    </div>";
    function trimLen(str) {
        if (str.length > 80) {
            str = str.slice(0, 77) + '...';
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
    function showError(error) {
        document.querySelector('pre#error-box')
            .innerHTML = quoteattr(error.toString());
    }
    function render(request, response) {
        var insertionPoint = document.querySelector('div.input-block');
        var html = '';
        var lastResp = '';
        request.commands.forEach(function (reqLine, i) {
            var resp = null;
            if (i < response.responses.length) {
                resp = response.responses[i];
            }
            console.log(response, resp);
            var error = !!resp.issue;
            var respLine = error ? resp.issue || 'err' : resp.replOutput || 'ok';
            lastResp = respLine;
            var preview = trimLen(respLine);
            html += "<div class=\"statement-block\">\n          &gt; <input class=\"repl-stmt\" value=\"".concat(quoteattr(reqLine), "\"><br>\n          <code class=\"").concat(error ? 'repl-out-error' : 'repl-out-ok', "\">").concat(quoteattr(preview), "</code>\n        </div>");
        });
        html += newInput;
        insertionPoint.innerHTML = html;
        document.querySelector('code#result').innerHTML = quoteattr(lastResp);
        readyForInput();
    }
    function submit() {
        var inputs = document.querySelectorAll('input.repl-stmt');
        var req = { 'commands': [] };
        inputs.forEach(function (el) {
            var inp = el;
            if (inp.value != null && inp.value.trim() != '') {
                req.commands.push(inp.value);
            }
        });
        var payload = JSON.stringify(req);
        console.log('submitting: %s', payload);
        window
            .fetch(endpoint, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: payload
        })
            .then(function (response) { return response.json(); })
            .then(function (data) { return render(req, data); })
            .catch(function (ex) { return showError(ex); });
    }
    function readyForInput() {
        var mainInput = document.querySelector('input.repl-stmt-new');
        mainInput.focus();
        mainInput.addEventListener('keyup', function (e) {
            if (e.code == 'Enter' && !e.metaKey && !e.ctrlKey) {
                e.preventDefault();
                submit();
            }
        });
    }
    window.addEventListener('load', function () {
        var evalBtn = document.querySelector('button#evaluate');
        var addBtn = document.querySelector('button#add-statement');
        if (!evalBtn || !addBtn) {
            console.log('Unable to bind control callbacks.');
            return;
        }
        var insertionPoint = document.querySelector('div.input-block');
        insertionPoint.innerHTML = newInput;
        evalBtn.addEventListener('click', function () { return submit(); });
        readyForInput();
    });
})();
