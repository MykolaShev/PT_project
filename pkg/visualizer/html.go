package visualizer

import (
	"encoding/json"
	"pt_project/pkg/ast"
	"pt_project/pkg/parser"
	"strings"
)

const htmlTemplate = `<!DOCTYPE html>
<html lang="uk">
<head>
    <meta charset="UTF-8">
    <title>SIPL Parse Tree vs AST Comparison</title>
    <script src="https://d3js.org/d3.v7.min.js"></script>
    <style>
        body { font-family: sans-serif; background-color: #f4f6f8; margin: 0; padding: 16px; }
        h2 { text-align: center; color: #222; margin: 0 0 6px 0; }
        p.hint { text-align: center; color: #555; font-size: 13px; margin: 0 0 14px 0; }
        .grid { display: grid; grid-template-columns: 1fr 1fr; gap: 16px; }
        .panel { background: #fff; border: 1px solid #d0d7de; border-radius: 8px; padding: 10px; box-shadow: 0 2px 6px rgba(0,0,0,0.06); }
        .panel h3 { margin: 0 0 8px 0; font-size: 15px; color: #1f2328; text-align: center; }
        svg { width: 100%; height: 78vh; border-top: 1px solid #eee; }
        .node circle { fill: #fff; stroke: #4682b4; stroke-width: 2.5px; cursor: pointer; transition: all 0.15s; }
        .node.terminal circle { stroke: #2e8b57; fill: #f0fff4; }
        .node.highlighted circle { fill: #ffedd5; stroke: #ea580c; stroke-width: 4px; r: 10; }
        .node text.label { font: 12px monospace; fill: #111; font-weight: bold; }
        .node text.span { font: 10px monospace; fill: #666; }
        .link { fill: none; stroke: #cbd5e1; stroke-width: 1.8px; }
    </style>
</head>
<body>
    <h2>SIPL Structural Comparison: Parse Tree (CST) &harr; Abstract Syntax Tree (AST)</h2>
    <p class="hint">Наведіть курсор на будь-який вузол, щоб підсвітити відповідний фрагмент коду [line:col] в обох деревах</p>
    <div class="grid">
        <div class="panel">
            <h3>Concrete Syntax Tree (Parse Tree)</h3>
            <svg id="cstSvg"></svg>
        </div>
        <div class="panel">
            <h3>Abstract Syntax Tree (AST)</h3>
            <svg id="astSvg"></svg>
        </div>
    </div>

    <script>
        const rawCst = {{CST_JSON}};
        const rawAst = {{AST_JSON}};

        function formatSpan(span) {
            if (!span || !span.start || !span.start.line) return "";
            return "[" + span.start.line + ":" + span.start.column + "]";
        }

        function mapCstNode(n) {
            if (!n) return { name: "nil", span: "" };
            const label = n.token ? n.symbol + " '" + n.token + "'" : n.symbol;
            return {
                name: label,
                span: formatSpan(n.span),
                isTerminal: Boolean(n.token),
                children: (n.children || []).map(mapCstNode)
            };
        }

        function mapAstNode(n) {
            if (!n) return { name: "nil", span: "" };
            const span = formatSpan(n.span);

            if (n.indentifier !== undefined) {
                return {
                    name: ":=",
                    span: span,
                    children: [
                        { name: n.indentifier.value || "id", span: formatSpan(n.indentifier.span), isTerminal: true },
                        mapAstNode(n.value)
                    ]
                };
            }
            if (n.left !== undefined && n.right !== undefined && n.operator !== undefined) {
                return { name: n.operator, span: span, children: [ mapAstNode(n.left), mapAstNode(n.right) ] };
            }
            if (n.left !== undefined && n.right !== undefined) {
                return { name: ";", span: span, children: [ mapAstNode(n.left), mapAstNode(n.right) ] };
            }
            if (n.then_block !== undefined) {
                let children = [
                    { name: "cond", span: span, children: [mapAstNode(n.condition)] },
                    { name: "then", span: span, children: [mapAstNode(n.then_block)] }
                ];
                if (n.else_block) {
                    children.push({ name: "else", span: span, children: [mapAstNode(n.else_block)] });
                }
                return { name: "if", span: span, children: children };
            }
            if (n.condition !== undefined && n.body !== undefined) {
                return {
                    name: "while",
                    span: span,
                    children: [
                        { name: "cond", span: span, children: [mapAstNode(n.condition)] },
                        { name: "body", span: span, children: [mapAstNode(n.body)] }
                    ]
                };
            }
            if (n.body !== undefined) {
                return { name: "begin...end", span: span, children: [ mapAstNode(n.body) ] };
            }
            if (n.name !== undefined) {
                return { name: "id(" + n.name + ")", span: span, isTerminal: true };
            }
            if (n.value !== undefined && typeof n.value === 'number') {
                return { name: "int(" + n.value + ")", span: span, isTerminal: true };
            }
            if (n.span !== undefined) {
                return { name: "skip", span: span, isTerminal: true };
            }
            return { name: "node", span: span };
        }

        function renderTree(svgSelector, rootData) {
            const width = 900;
            const height = 650;
            const svg = d3.select(svgSelector)
                .attr("viewBox", "0 0 " + width + " " + height)
                .append("g")
                .attr("transform", "translate(40,40)");

            const root = d3.hierarchy(rootData);
            const treeLayout = d3.tree().size([width - 80, height - 90]);
            treeLayout(root);

            svg.selectAll('.link')
                .data(root.links())
                .enter().append('path')
                .attr('class', 'link')
                .attr('d', d3.linkVertical().x(d => d.x).y(d => d.y));

            const node = svg.selectAll('.node')
                .data(root.descendants())
                .enter().append('g')
                .attr('class', d => 'node' + (d.data.isTerminal ? ' terminal' : ''))
                .attr('data-span', d => d.data.span)
                .attr('transform', d => 'translate(' + d.x + ',' + d.y + ')')
                .on('mouseenter', function(event, d) {
                    if (!d.data.span) return;
                    d3.selectAll('.node').classed('highlighted', n => n.data.span === d.data.span);
                })
                .on('mouseleave', function() {
                    d3.selectAll('.node').classed('highlighted', false);
                });

            node.append('circle').attr('r', 7);
            node.append('title').text(d => d.data.name + " " + d.data.span);

            node.append('text')
                .attr('class', 'label')
                .attr("dy", d => d.children && d.children.length ? -12 : 18)
                .style("text-anchor", "middle")
                .text(d => d.data.name);

            node.append('text')
                .attr('class', 'span')
                .attr("dy", d => d.children && d.children.length ? -24 : 30)
                .style("text-anchor", "middle")
                .text(d => d.data.span);
        }

        if (rawCst) {
            renderTree("#cstSvg", mapCstNode(rawCst));
        }
        if (rawAst && rawAst.statements) {
            renderTree("#astSvg", {
                name: "Program",
                span: rawAst.statements.length ? formatSpan(rawAst.statements[0].span) : "",
                children: rawAst.statements.map(mapAstNode)
            });
        }
    </script>
</body>
</html>`

// GenerateHTML injects both Parse Tree (CST) and AST JSON into the D3.js comparison view.
func GenerateHTML(program *ast.Program, cst *parser.ParseNode) string {
	cstJSON, _ := json.Marshal(cst)
	astJSON, _ := json.Marshal(program)
	out := strings.Replace(htmlTemplate, "{{CST_JSON}}", string(cstJSON), 1)
	return strings.Replace(out, "{{AST_JSON}}", string(astJSON), 1)
}