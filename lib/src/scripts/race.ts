import * as d3 from "d3";

export class Race {
  private readonly factor: number = 40;
  private width: number;
  private height: number;
  private data: any;

  constructor(data: any) {
    this.width = data.track.width;
    this.height = data.track.height;
    this.data = data;
  }

  create(): SVGSVGElement {
    const svg = d3
      .create("svg")
      .attr("width", this.width * this.factor)
      .attr("height", this.height * this.factor);

    this.appendGrid(svg);
    this.appendTrack(svg);

    return svg.node()!;
  }

  private appendPlayers(
    svg: d3.Selection<SVGSVGElement, undefined, null, undefined>,
    data: any,
  ) {
    data.players.forEach((p: any) => {
      svg
        .append("circle")
        .attr("cx", p.path[0].x * this.factor)
        .attr("cy", p.path[0].y * this.factor)
        .attr("r", 10)
        .attr("fill", p.color);
    });
  }

  private appendTrack(
    svg: d3.Selection<SVGSVGElement, undefined, null, undefined>,
  ) {
    const outer = [
      ...this.data.track.outer.map((p: any) => [
        p.x * this.factor,
        p.y * this.factor,
      ]),
      [
        this.data.track.outer[0].x * this.factor,
        this.data.track.outer[0].y * this.factor,
      ],
    ];
    const inner = [
      ...this.data.track.inner.map((p: any) => [
        p.x * this.factor,
        p.y * this.factor,
      ]),
      [
        this.data.track.inner[0].x * this.factor,
        this.data.track.inner[0].y * this.factor,
      ],
    ];
    const finish = this.data.track.finish.map((p: any) => [
      p.x * this.factor,
      p.y * this.factor,
    ]);
    const line = d3.line();

    svg
      .append("path")
      .attr(
        "d",
        line([
          [0, 0],
          [this.width * this.factor, 0],
          [this.width * this.factor, this.height * this.factor],
          [0, this.height * this.factor],
          [0, 0],
        ]),
      )
      .attr("stroke", "black")
      .attr("fill", "none")
      .attr("stroke-width", 4);

    svg
      .append("path")
      .attr("d", line(finish))
      .attr("stroke", "black")
      .attr("fill", "none")
      .attr("stroke-width", 12)
      .style("stroke-dasharray", "20, 10");

    svg
      .append("path")
      .attr("d", line(outer))
      .attr("stroke", "black")
      .attr("fill", "none")
      .attr("stroke-width", 8);

    svg
      .append("path")
      .attr("d", line(inner))
      .attr("stroke", "black")
      .attr("fill", "none")
      .attr("stroke-width", 8);
  }

  private appendGrid(
    svg: d3.Selection<SVGSVGElement, undefined, null, undefined>,
  ) {
    const line = d3.line();

    for (let x = 1; x < this.data.track.width; x++) {
      svg
        .append("path")
        .attr(
          "d",
          line([
            [x * this.factor, 0],
            [x * this.factor, this.height * this.factor],
          ]),
        )
        .attr("stroke", "lightgray")
        .attr("fill", "none");
    }

    for (let y = 1; y < this.data.track.height; y++) {
      svg
        .append("path")
        .attr(
          "d",
          line([
            [0, y * this.factor],
            [this.width * this.factor, y * this.factor],
          ]),
        )
        .attr("stroke", "lightgray")
        .attr("fill", "none");
    }
  }
}
