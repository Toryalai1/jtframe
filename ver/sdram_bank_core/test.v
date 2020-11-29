`timescale 1ns / 1ps

module test;

reg        rst, clk, init_done, waiting;

reg [22:0] addr_req;
reg        rd_req;
reg [ 1:0] ba_req;

wire [31:0] dout;
wire        ack, rdy;
wire [ 1:0] ba_rdy;

// sdram pins
wire [15:0] sdram_dq;
wire [12:0] sdram_a;
wire [ 1:0] sdram_dqm;
wire [ 1:0] sdram_ba;
wire        sdram_nwe;
wire        sdram_ncas;
wire        sdram_nras;
wire        sdram_ncs;
wire        sdram_cke;

always @(posedge clk, posedge rst) begin
    if( rst ) begin
        addr_req <= 23'd0;
        rd_req   <= 0;
        ba_req   <= 2'd0;
        waiting  <= 0;
    end else if(init_done) begin
        if( !waiting ) begin
            addr_req[19:0] <= $urandom;
            ba_req   <= $urandom;
            rd_req   <= 1;
            waiting  <= 1;
        end else if( ack ) begin
            waiting <= 0;
            rd_req <= 0;
        end
    end
end

jtframe_sdram_bank_core uut(
    .rst        ( rst           ),
    .clk        ( clk           ),
    .addr       ( addr_req      ),
    .rd         ( rd_req        ),
    .wr         ( 1'b0          ),
    .ba_rq      ( ba_req        ),
    .ack        ( ack           ),
    .rdy        ( rdy           ),
    .ba_rdy     ( ba_rdy        ),
    .din        ( 16'd0         ),
    .dout       ( dout          ),
    // SDRAM pins
    .sdram_dq   ( sdram_dq      ),
    .sdram_a    ( sdram_a       ),
    .sdram_dqml ( sdram_dqm[0]  ),
    .sdram_dqmh ( sdram_dqm[1]  ),
    .sdram_ba   ( sdram_ba      ),
    .sdram_nwe  ( sdram_nwe     ),
    .sdram_ncas ( sdram_ncas    ),
    .sdram_nras ( sdram_nras    ),
    .sdram_ncs  ( sdram_ncs     ),
    .sdram_cke  ( sdram_cke     )
);


mt48lc16m16a2 sdram(
    .Clk        ( clk       ),
    .Cke        ( sdram_cke ),
    .Dq         ( sdram_dq  ),
    .Addr       ( sdram_a   ),
    .Ba         ( sdram_ba  ),
    .Cs_n       ( sdram_ncs ),
    .Ras_n      ( sdram_nras),
    .Cas_n      ( sdram_ncas),
    .We_n       ( sdram_nwe ),
    .Dqm        ( sdram_dqm ),
    .downloading( 1'b0      ),
    .VS         ( 1'b0      ),
    .frame_cnt  ( 0         )
);

initial begin
    clk=0;
    forever #5.2 clk=~clk;
end

initial begin
    init_done = 0;
    #104_500 init_done = 1;
end

initial begin
    rst=1;
    #100 rst=0;
    #500_000 $finish;
end

initial begin
    $dumpfile("test.lxt");
    $dumpvars;
end

endmodule