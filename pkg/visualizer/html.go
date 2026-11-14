package visualizer

import (
	"encoding/json"
	"pt_project/pkg/ast"
	"strings"
)

const htmlTemplate = `<!DOCTYPE html>
<html lang="uk">
<head>
    <meta charset="UTF-8">
    <title>SIPL AST Visualizer</title>
    <script src="https://d3js.org/d3.v7.min.js"></script>
    <style>
        body { font-family: sans-serif; display: flex; flex-direction: column; align-items: center; background-color: #f9f9f9; margin: 0; padding: 20px; }
        h2 { color: #333; margin-bottom: 10px; }
        svg { border: 1px solid #ccc; background-color: #fff; width: 95vw; height: 80vh; border-radius: 8px; box-shadow: 0 4px 6px rgba(0,0,0,0.1); }
        .node circle { fill: #fff; stroke: #4682b4; stroke-width: 3px; }
        .node text { font: 13px monospace; fill: #333; font-weight: bold; }
        .link { fill: none; stroke: #ccc; stroke-width: 2px; }
    </style>
</head>
<body>
    <h2>SIPL Interactive AST Visualizer</h2>
    <svg id="treeSvg"></svg>

    <script>
        const rawJson = {{AST_JSON}};

        // NOTE: the Go side marshals struct fields using their lower-case
        // "json:" tags (see pkg/ast/ast.go), not the Go field names, so the
        // checks below must match those tag names exactly ("indentifier",
        // "then_block", "else_block", "statements", etc.) rather than the
        // capitalized Go identifiers the previous version of this file used.
        function mapNode(n) {
            if (!n) return { name: "nil" };

            // AssignStmt: { indentifier: {value}, value: Expression }
            if (n.indentifier !== undefined) {
                return {
                    name: ":=",
                    children: [ { name: n.indentifier.value || "id" }, mapNode(n.value) ]
                };
            }
            // BinaryExpr: { left, operator, right } (checked before SequenceStmt,
            // which also has left/right but no operator)
            if (n.left !== undefined && n.right !== undefined && n.operator !== undefined) {
                return {
                    name: n.operator,
                    children: [ mapNode(n.left), mapNode(n.right) ]
                };
            }
            // SequenceStmt: { left, right } (Statements, no operator)
            if (n.left !== undefined && n.right !== undefined) {
                return {
                    name: ";",
                    children: [ mapNode(n.left), mapNode(n.right) ]
                };
            }
            // IfStmt: { condition, then_block, else_block }
            if (n.then_block !== undefined) {
                let children = [ {name: "cond", children: [mapNode(n.condition)]}, {name: "then", children: [mapNode(n.then_block)]} ];
                if (n.else_block) {
                    children.push({name: "else", children: [mapNode(n.else_block)]});
                }
                return { name: "if", children: children };
            }
            // WhileStmt: { condition, body }
            if (n.condition !== undefined && n.body !== undefined) {
                return {
                    name: "while",
                    children: [ {name: "cond", children: [mapNode(n.condition)]}, {name: "body", children: [mapNode(n.body)]} ]
                };
            }
            // BlockStmt (begin...end): { body } only, no condition
            if (n.body !== undefined) {
                return { name: "begin...end", children: [ mapNode(n.body) ] };
            }
            // IdentifierExpr: { name }
            if (n.name !== undefined) {
                return { name: "id(" + n.name + ")" };
            }
            // NumberExpr: { value } (number)
            if (n.value !== undefined && typeof n.value === 'number') {
                return { name: "int(" + n.value + ")" };
            }
            // SkipStmt: only { span }, nothing else distinguishing
            if (n.span !== undefined) {
                return { name: "skip" };
            }

            return { name: "node" };
        }

        function drawTree() {
            if (!rawJson || !rawJson.statements) {
                console.error("Invalid AST JSON structure", rawJson);
                return;
            }

            let rootData = { 
                name: "Program", 
                children: rawJson.statements.map(mapNode) 
            };
            
            const width = window.innerWidth * 0.9;
            const height = window.innerHeight * 0.75;
            
            const svg = d3.select("#treeSvg")
                .attr("viewBox", '0 0 ' + width + ' ' + height)
                .append("g")
                .attr("transform", "translate(40,40)");

            const root = d3.hierarchy(rootData);
            const treeLayout = d3.tree().size([width - 80, height - 80]);
            treeLayout(root);

            // Draw links
            svg.selectAll('.link')
                .data(root.links())
                .enter().append('path')
                .attr('class', 'link')
                .attr('d', d3.linkVertical().x(d => d.x).y(d => d.y));

            // Draw nodes
            const node = svg.selectAll('.node')
                .data(root.descendants())
                .enter().append('g')
                .attr('class', 'node')
                .attr('transform', d => 'translate(' + d.x + ',' + d.y + ')');

            node.append('circle').attr('r', 8);
            
            node.append('text')
                .attr("dy", d => d.children ? -15 : 22)
                .style("text-anchor", "middle")
                .text(d => d.data.name);
        }

        drawTree();
    </script>
</body>
</html>`

// GenerateHTML injects the AST JSON into the D3.js HTML template.
func GenerateHTML(program *ast.Program) string {
	astJSON, _ := json.Marshal(program)
	return strings.Replace(htmlTemplate, "{{AST_JSON}}", string(astJSON), 1)
}