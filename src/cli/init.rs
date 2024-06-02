use clap::{Parser, Subcommand};

#[derive(Parser, Debug)]
struct Args {
    file: Option<String>,

    #[command(subcommand)]
    command: Option<Commands>,
}

#[derive(Subcommand, Debug)]
enum Commands {
    /// Manage the editor settings.
    Config {
        /// View the editor settings.
        #[clap(short, long)]
        view: Option<String>,

        /// Edit the editor settings.
        #[clap(short, long)]
        edit: Option<String>,
    },

    /// View or downlaod the documentation.
    Docs {
        /// View the documentation.
        #[clap(short, long)]
        view: Option<String>,

        /// Download the documentation.
        #[clap(short, long)]
        download: Option<String>,
    },
}

pub fn cli() {
    let args = Args::parse();

    match args.file {
        Some(file) => {
            if !file.is_empty() {
                // TODO: Open the file.
            }
        }
        None => {
            // TODO: Open a new file in blank.
        }
    }

    match args.command {
        Some(Commands::Config { view, edit }) => {
            if let Some(_viewww) = view {
                // TODO: View the editor settings.
            } else if let Some(_edit) = edit {
                // TODO: Edit the editor settings.
            }
        }
        Some(Commands::Docs { view, download }) => {
            if let Some(_view) = view {
                // TODO: View the documentation.
            } else if let Some(_download) = download {
                // TODO: Download the documentation.
            }
        }
        None => {}
    }
}
