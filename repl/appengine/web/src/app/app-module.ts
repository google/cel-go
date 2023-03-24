import { NgModule } from '@angular/core';
import { BrowserModule } from '@angular/platform-browser';
import { BrowserAnimationsModule } from '@angular/platform-browser/animations';
import { HttpClientModule } from '@angular/common/http';
import { MatFormFieldModule } from '@angular/material/form-field';
import { MatInputModule } from '@angular/material/input';
import { MatSidenavModule } from '@angular/material/sidenav';
import { MatButtonModule } from '@angular/material/button';
import { AppComponent } from './app-component';
import { ReplConsoleModule } from './repl_console/repl-console-module';
import { ReferencePanelModule } from './reference_panel/reference-panel-module';
import { SharedModule } from './shared/shared-module';

@NgModule({
  declarations: [
    AppComponent,
  ],
  imports: [
    BrowserModule,
    HttpClientModule,
    SharedModule,
    MatFormFieldModule,
    ReplConsoleModule,
    BrowserAnimationsModule,
    MatInputModule,
    MatSidenavModule,
    MatButtonModule,
    ReferencePanelModule
  ],
  providers: [],
  bootstrap: [AppComponent]
})
export class AppModule { }
