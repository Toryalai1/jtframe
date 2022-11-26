/*  This file is part of JTFRAME.
    JTFRAME program is free software: you can redistribute it and/or modify
    it under the terms of the GNU General Public License as published by
    the Free Software Foundation, either version 3 of the License, or
    (at your option) any later version.

    JTFRAME program is distributed in the hope that it will be useful,
    but WITHOUT ANY WARRANTY; without even the implied warranty of
    MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
    GNU General Public License for more details.

    You should have received a copy of the GNU General Public License
    along with JTFRAME.  If not, see <http://www.gnu.org/licenses/>.

    Author: Jose Tejada Gomez. Twitter: @topapate
    Version: 1.0
    Date: 30-10-2022 */

// Frame buffer built on top of two line buffers
// the frame is stored in two PSRAM chips

// This module is not fully tested yet
module jtframe_lfbuf_ddr #(parameter
    DW      =  16,
    VW      =   8,
    HW      =   9
)(
    input               rst,     // hold in reset for >150 us
    input               clk,
    input               pxl_cen,

    // video status
    input      [VW-1:0] vrender,
    input      [HW-1:0] hdump,
    input               vs,
    input               lhbl,
    input               lvbl,

    // core interface
    input      [HW-1:0] ln_addr,
    input      [DW-1:0] ln_data,
    input               ln_done,
    output              ln_hs,
    output     [DW-1:0] ln_pxl,
    output     [VW-1:0] ln_v,
    input               ln_we,

    // DDR3 RAM
    output              ddram_clk,
    input               ddram_busy,
    output      [7:0]   ddram_burstcnt,
    output     [28:0]   ddram_addr,
    input      [63:0]   ddram_dout,
    input               ddram_dout_ready,
    output              ddram_rd,
    output     [63:0]   ddram_din,
    output      [7:0]   ddram_be,
    output              ddram_we,
    // Status
    input       [7:0]   st_addr,
    output      [7:0]   st_dout
);

wire          frame, fb_clr, fb_done, line, scr_we;
wire [HW-1:0] fb_addr, rd_addr;
wire [  15:0] fb_din, fb_dout;

jtframe_lfbuf_ddr_ctrl #(.HW(HW),.VW(VW)) u_ctrl (
    .rst        ( rst       ),
    .clk        ( clk       ),
    .pxl_cen    ( pxl_cen   ),

    .lhbl       ( lhbl      ),
    .ln_done    ( ln_done   ),
    .vrender    ( vrender   ),
    .ln_v       ( ln_v      ),
    // data written to external memory
    .frame      ( frame     ),
    .fb_addr    ( fb_addr   ),
    .rd_addr    ( rd_addr   ),
    .fb_din     ( fb_din    ),
    .fb_dout    ( fb_dout   ),
    .fb_clr     ( fb_clr    ),
    .fb_done    ( fb_done   ),

    // data read from external memory to screen buffer
    // during h blank
    .line       ( line      ),
    .scr_we     ( scr_we    ),

    .ddram_clk  ( ddram_clk     ),
    .ddram_busy ( ddram_busy    ),
    .ddram_addr ( ddram_addr    ),
    .ddram_dout ( ddram_dout    ),
    .ddram_rd   ( ddram_rd      ),
    .ddram_din  ( ddram_din     ),
    .ddram_be   ( ddram_be      ),
    .ddram_we   ( ddram_we      ),
    .ddram_burstcnt  ( ddram_burstcnt    ),
    .ddram_dout_ready( ddram_dout_ready  ),
    .st_addr    ( st_addr       ),
    .st_dout    ( st_dout       )
);

jtframe_lfbuf_line #(.DW(DW),.HW(HW),.VW(VW)) u_line(
    .rst        ( rst       ),
    .clk        ( clk       ),
    // video status
    .vrender    ( vrender   ),
    .hdump      ( hdump     ),
    .vs         ( vs        ),   // vertical sync, the buffer is swapped here
    .lvbl       ( lvbl      ),   // vertical blank, active low

    // core interface
    .ln_hs      ( ln_hs     ),
    .ln_v       ( ln_v      ),
    .ln_addr    ( ln_addr   ),
    .ln_data    ( ln_data   ),
    .ln_we      ( ln_we     ),
    .ln_pxl     ( ln_pxl    ),

    // data written to external memory
    .frame      ( frame     ),
    .fb_addr    ( fb_addr   ),
    .rd_addr    ( rd_addr   ),
    .fb_din     ( fb_din    ),
    .fb_dout    ( fb_dout   ),
    .fb_clr     ( fb_clr    ),
    .fb_done    ( fb_done   ),

    // data read from external memory to screen buffer
    // during h blank
    .line       ( line      ),
    .scr_we     ( scr_we    )
);

endmodule